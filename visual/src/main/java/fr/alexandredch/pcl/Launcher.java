package fr.alexandredch.pcl;

import org.apache.commons.codec.DecoderException;
import org.apache.commons.codec.binary.Hex;
import org.w3c.dom.Document;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;
import visual.EmulatorLogFile;
import visual.HeadlessController;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Stream;

public class Launcher {

    /* size of the print buffer */
    private static final int instMemSize = 0x10000;
    /* we make sure that the output buffer is always the first symbol in memory */
    private static final int outputBufferAddress = instMemSize;
    /* VisUAL offsets line numbers by one for some reason */
    private static final List<Integer> breakpoints = new ArrayList<>(); //= List.of(13 - 1);

    /* array of all word addresses in the output buffer */
    private static String[] getOutputRange() {
        return Stream.iterate(outputBufferAddress, n -> n + 4)
                .limit(0x1000 / 4)
                .map(n -> String.format("0x%X", n))
                .toArray(String[]::new);
    }

    public static void executeAndParseOutput(String assemblyFile) {
        // Search for the line containing "BNE     PRINTLN_LOOP"
        File file = new File(assemblyFile);
        if (!file.exists()) {
            throw new IllegalArgumentException("File does not exist");
        }
        try {
            List<String> lines = java.nio.file.Files.readAllLines(file.toPath());
            for (int i = 0; i < lines.size(); i++) {
                if (lines.get(i).contains("STRB    R2, [R1, #-1]")) {
                    breakpoints.add(i);
                }
            }
        } catch (IOException e) {
            throw new RuntimeException(e);
        }

        EmulatorLogFile.configureLogging("", true, false, false, false, false, false,
                true, false, getOutputRange());
        HeadlessController.setLogMode(EmulatorLogFile.LogMode.BREAKPOINT);
        HeadlessController.setBreakpoints(breakpoints); // these are the lines where the program will stop
        HeadlessController.setInstMemSize(instMemSize);
        String logFile = String.format("%s_log.xml", assemblyFile);

        // HeadlessController#runFile will try to exit once it is done, there is no good way to prevent this
        // A workaround is adding a shutdown hook to parse the output
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            List<String> output = parseOutput(logFile);
            System.out.println("---- PROGRAM OUTPUT ----");
            output.forEach(System.out::print);
            System.out.println("\n---- END PROGRAM OUTPUT ----");
        }));

        HeadlessController.runFile(assemblyFile, logFile);
    }

    public static List<String> parseOutput(String XMLPath) {
        List<String> outputs = new ArrayList<>();
        try {
            File XMLFile = new File(XMLPath);
            DocumentBuilderFactory dbFactory = DocumentBuilderFactory.newInstance();
            DocumentBuilder dBuilder = dbFactory.newDocumentBuilder();
            Document doc = dBuilder.parse(XMLFile);
            doc.getDocumentElement().normalize();
            NodeList lines = doc.getElementsByTagName("line");
            for (int i = 0; i < lines.getLength(); i++) {
                System.out.println("Parsing line " + i);
                Node line = lines.item(i);
                outputs.add(reverseString(parseLine(line)));
            }
            System.out.println("Parsed " + lines.getLength() + " lines");
        } catch (Exception e) {
            e.printStackTrace();
        }
        return outputs;
    }

    private static String parseLine(Node line) {
        NodeList children = line.getChildNodes();
        ByteArrayOutputStream bytes = new ByteArrayOutputStream();
        for (int j = 0; j < children.getLength(); j++) {
            Node child = children.item(j);
            if (child.getNodeName().equals("word")) {
                try {
                    // The 'A' at the end represents a new line
                    if (child.getTextContent().replaceAll("0x", "").equals("A")) {
                        // We reached the end, the string is empty
                        // Return a new line
                        bytes.write(0x0A);
                    } else {
                        if (child.getTextContent().replaceAll("0x", "").length() % 2 == 1) {
                            // The decode needs to have a length that is a multiple of 2, so we should have an 'A'
                            // at the end if it is not the case
                            // Any other case is just avoided
                            if (child.getTextContent().replaceAll("0x", "").charAt(0) == 'A') {
                                String hex = child.getTextContent().replaceAll("0x", "").substring(1);
                                byte[] toBytes = Hex.decodeHex(hex.replaceAll("0x", "").toCharArray());
                                toBytes = reverseByteArray(toBytes);
                                bytes.write(toBytes);
                                bytes.write(0x0A);
                            }
                        } else {
                            byte[] toBytes = Hex.decodeHex(child.getTextContent().replaceAll("0x", "").toCharArray());
                            toBytes = reverseByteArray(toBytes);
                            bytes.write(toBytes);
                        }
                    }
                } catch (DecoderException | IOException e) {
                    throw new RuntimeException(e);
                }
            }
        }
        return bytes.toString(StandardCharsets.UTF_8);
    }

    private static byte[] reverseByteArray(byte[] array) {
        byte[] reversed = new byte[array.length];
        for (int i = 0; i < array.length; i++) {
            reversed[i] = array[array.length - i - 1];
        }
        return reversed;
    }

    private static String reverseString(String s) {
        if (s.charAt(s.length() - 1) == '\n') {
            // Keep the line break at the end
            return new StringBuilder(s.substring(0, s.length() - 1)).reverse().append('\n').toString();
        }
        return new StringBuilder(s).reverse().toString();
    }

    public static void run(String assemblyFile) {
        System.out.println("---- RUNNING PROGRAM ----");
        executeAndParseOutput(assemblyFile);
    }
}
