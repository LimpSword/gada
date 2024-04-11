plugins {
    id("java")
    id ("com.github.johnrengelman.shadow") version "8.1.1"
}

group = "fr.alexandredch"
version = "1.0-SNAPSHOT"

repositories {
    mavenCentral()
}

configurations {
    compileOnly {
        extendsFrom(configurations.annotationProcessor.get())
    }
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(11))
    }
}

dependencies {
    implementation(files("libs/visual.jar"))
    implementation("commons-codec:commons-codec:1.16.1")
}

tasks.withType<Jar> {
    duplicatesStrategy = DuplicatesStrategy.EXCLUDE
    configurations["compileClasspath"].forEach { file: File ->
        from(zipTree(file.absoluteFile))
    }

    manifest {
        attributes["Main-Class"] = "fr.alexandredch.pcl.Application"
    }

    // Set JAR file name
    archiveFileName.set("pcl.jar")
}

tasks.test {
    useJUnitPlatform()
}