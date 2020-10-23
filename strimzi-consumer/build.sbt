lazy val AkkaVersion = "2.6.9"
lazy val AkkaManagementVersion = "1.0.8"
lazy val AlpakkaKafkaVersion = "2.0.5"
lazy val LogbackVersion = "1.2.3"

lazy val commonScalacOptions = Seq(
  "-encoding", "UTF-8",
  "-deprecation",
  "-feature",
  "-unchecked",
  "-Xlint",
  "-Ywarn-unused:imports",
  "-Ywarn-dead-code"
)

lazy val commonJavacOptions = Seq(
  "-Xlint:unchecked",
  "-Xlint:deprecation"
)

lazy val commonSettings = Seq(
  Compile / scalacOptions ++= commonScalacOptions,
  Compile / javacOptions ++= commonJavacOptions,
  run / javaOptions ++= Seq("-Xms128m", "-Xmx1024m"),
  run / fork := false,
  Global / cancelable := false
)
// https://www.scala-sbt.org/sbt-native-packager/formats/docker.html
// https://github.com/sbt/sbt-native-packager
// dockerAlias defaults to: [dockerRepository/][dockerUsername/][packageName]:[version]
lazy val dockerSettings = Seq(
  dockerBaseImage := "openjdk:11-jre-slim",
  packageName in Docker := "strimzi-consumer",
  maintainer in Docker := "Sina Nourian <sina.nourian@gmail.com>",
  version in Docker := "v1",
  dockerUsername := Some("nrsina")
  //dockerAlias := DockerAlias(None, dockerUsername.value, (packageName in Docker).value, Some((version in Docker).value))
)

lazy val producer = (project in file("."))
  .settings(commonSettings)
  .settings(PB.targets in Compile := Seq(scalapb.gen(flatPackage = true) -> (sourceManaged in Compile).value))
  .enablePlugins(JavaAppPackaging)
  .settings(dockerSettings)
  .settings(
    name := "strimzi-consumer",
    version := "0.1",
    scalaVersion := "2.13.3",
    mainClass in(Compile, run) := Some("com.snourian.strimzi.consumer.Application"),
    libraryDependencies ++= Seq(
      "com.typesafe.akka" %% "akka-actor-typed" % AkkaVersion,
      "com.typesafe.akka" %% "akka-stream" % AkkaVersion,
      "com.typesafe.akka" %% "akka-stream-kafka" % AlpakkaKafkaVersion,
      "com.typesafe.akka" %% "akka-slf4j" % AkkaVersion,
      "com.lightbend.akka.management" %% "akka-management" % AkkaManagementVersion,
      "com.thesamet.scalapb" %% "scalapb-runtime" % scalapb.compiler.Version.scalapbVersion % "protobuf",
      "ch.qos.logback" % "logback-classic" % LogbackVersion)
  )