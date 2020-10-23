package com.snourian.strimzi.consumer

import akka.actor.ActorSystem
import akka.kafka.ConsumerSettings
import org.apache.kafka.clients.consumer.ConsumerConfig
import org.apache.kafka.common.serialization.{ByteArrayDeserializer, StringDeserializer}

import scala.concurrent.duration.DurationInt

object ConsumerProperties {

  def apply(system: ActorSystem): ConsumerProperties = {
    val config = system.settings.config.getConfig("strimzi-consumer")
    ConsumerProperties(
      config.getString("bootstrap-servers"),
      config.getString("group-id"),
      config.getString("enable-auto-commit"),
      config.getString("auto-offset-reset"),
      config.getString("topic"),
      config.getInt("parallelism"),
      config.getLong("sleep-time-ms"),
      system
    )
  }
}

final case class ConsumerProperties(bootstrapServers: String, groupId: String, enableAutoCommit: String,
                                    autoOffsetReset: String, topic: String, parallelism: Int,
                                    sleepTime: Long, system: ActorSystem) {

  def toConsumerSettings: ConsumerSettings[String, Array[Byte]] = {
    ConsumerSettings(system, new StringDeserializer, new ByteArrayDeserializer)
      .withBootstrapServers(bootstrapServers)
      .withGroupId(groupId)
      .withProperty(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, autoOffsetReset)
      .withProperty(ConsumerConfig.ENABLE_AUTO_COMMIT_CONFIG, enableAutoCommit)
      .withStopTimeout(2.seconds)
  }
}