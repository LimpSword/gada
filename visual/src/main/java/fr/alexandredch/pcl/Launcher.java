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
    private static final List<Integer> breakpoints = List.of(13 - 1);

    /* array of all word addresses in the output buffer */
    private static String[] getOutputRange() {
        return Stream.iterate(outputBufferAddress, n -> n + 4)
                .limit(0x1000 / 4)
                .map(n -> String.format("0x%X", n))
                .toArray(String[]::new);
    }

    public static void executeAndParseOutput(String assemblyFile) {
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
            output.forEach(System.out::println);
            System.out.println("---- END PROGRAM OUTPUT ----");
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
                outputs.add(parseLine(line));
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
                    // To be fair, it would have been a better idea not to add the 'A' at the end of the print buffer
                    if (child.getTextContent().replaceAll("0x", "").equals("A")) {
                        // We reached the end, the string is empty
                        break;
                    }
                    if (child.getTextContent().replaceAll("0x", "").length() % 2 == 1) {
                        // The decode needs to have a length that is a multiple of 2, so we should have an 'A'
                        // at the end if it is not the case
                        // Any other case is just avoided
                        if (child.getTextContent().replaceAll("0x", "").charAt(0) == 'A') {
                            String hex = child.getTextContent().replaceAll("0x", "").substring(1);
                            byte[] toBytes = Hex.decodeHex(hex.replaceAll("0x", "").toCharArray());
                            toBytes = reverseByteArray(toBytes);
                            System.out.println(child.getTextContent() + " -> " + new String(toBytes, StandardCharsets.UTF_8));
                            bytes.write(toBytes);
                        }
                        break;
                    }
                    byte[] toBytes = Hex.decodeHex(child.getTextContent().replaceAll("0x", "").toCharArray());
                    toBytes = reverseByteArray(toBytes);
                    System.out.println(child.getTextContent() + " -> " + new String(toBytes, StandardCharsets.UTF_8));
                    bytes.write(toBytes);
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

    public static void run(String assemblyFile) {
        System.out.println("---- RUNNING PROGRAM ----");
        executeAndParseOutput(assemblyFile);
    }
}
