package fr.alexandredch.pcl;

import org.w3c.dom.Document;
import org.w3c.dom.Node;
import org.w3c.dom.NodeList;
import visual.EmulatorLogFile;
import visual.HeadlessController;

import javax.xml.parsers.DocumentBuilder;
import javax.xml.parsers.DocumentBuilderFactory;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.stream.Stream;

public class Launcher {

    private static final int instMemSize = 0x10000;
    /* we make sure that the output buffer is always the first symbol in memory */
    private static final int outputBufferAddress = instMemSize;
    /* VisUAL offsets line numbers by one for some reason */
    private static final List<Integer> breakpoints = List.of(11 - 1, 25 - 1);

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
        HeadlessController.setBreakpoints(breakpoints);
        HeadlessController.setInstMemSize(instMemSize);
        String logFile = String.format("%s_log.xml", assemblyFile);

        // shutdown hook: parse output and exit
        Runtime.getRuntime().addShutdownHook(new Thread(() -> {
            List<String> output = parseOutput(logFile);
            System.out.println("---- PROGRAM OUTPUT ----");
            output.forEach(System.out::print);
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
                Node line = lines.item(i);
                outputs.add(parseLine(line));
            }
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
                boolean more = hexStringToByteArray(bytes, child.getTextContent());
                if (!more) break;
            }
        }
        return bytes.toString(StandardCharsets.UTF_8);
    }

    private static boolean hexStringToByteArray(ByteArrayOutputStream bytes, String hex) {
        String hexValue = hex.substring(2);
        int l = hexValue.length();
        if (l % 2 == 1) {
            /* pad to a pair number of hex digits */
            hexValue = "0" + hexValue;
            l++;
        }
        for (int i = l - 2; i >= 0; i -= 2) {
            int currentByte = (Character.digit(hexValue.charAt(i), 16) << 4)
                              + Character.digit(hexValue.charAt(i + 1), 16);
            if (currentByte == 0) return false;
            bytes.write(currentByte);
        }
        return true;
    }

    public static void run(String assemblyFile) {
        System.out.println("---- RUNNING PROGRAM ----");
        executeAndParseOutput(assemblyFile);
    }
}
