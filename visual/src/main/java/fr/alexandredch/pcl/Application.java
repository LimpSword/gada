package fr.alexandredch.pcl;

import java.io.File;

public final class Application {

    public static void main(String[] args) {
        if (args.length != 1) {
            String jarName = new File(Application.class.getProtectionDomain().getCodeSource().getLocation().getPath()).getName();

            System.err.printf("Usage: java -jar %s <assembly file>%n", jarName);
            System.exit(1);
        }

        File file = new File(args[0]);
        if (!file.exists()) {
            System.err.println("File not found: " + args[0]);
            System.exit(1);
        }
        Launcher.run(args[0]);
    }
}
