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
}

tasks.withType<Jar> {
    manifest {
        attributes["Main-Class"] = "fr.alexandredch.pcl.Application"
    }

    from (configurations.runtimeClasspath.get().filter { it.name.endsWith("jar") }.map { zipTree(it) })
}

tasks.test {
    useJUnitPlatform()
}