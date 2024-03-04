package fr.alexandredch.pcl;

import java.io.File;

public final class Application {

    public static void main(String[] args) {
        File file = new File(args[0]);
        if (!file.exists()) {
            System.err.println("File not found: " + args[0]);
            System.exit(1);
        }
        Launcher.run(args[0]);
    }
}
