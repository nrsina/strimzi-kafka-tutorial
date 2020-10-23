package com.snourian.strimzi.consumer

import akka.actor.ActorSystem
import akka.actor.typed.Behavior
import akka.actor.typed.scaladsl.Behaviors
import akka.actor.typed.scaladsl.adapter.TypedActorSystemOps
import akka.kafka.scaladsl.{Committer, Consumer}
import akka.kafka.{CommitterSettings, Subscriptions}
import akka.management.scaladsl.AkkaManagement
import com.snourian.strimzi.producer.consumer.MyMessage

import scala.concurrent.{ExecutionContextExecutor, Future}
import scala.util.Try

object ConsumerClient {

  sealed trait Command

  private final case class KafkaConsumerStopped(reason: Try[Any]) extends Command

  def apply(): Behavior[Nothing] = {
    Behaviors.setup[Command] { ctx =>
      ctx.log.debug("Initializing ConsumerClient actor...")
      // Alpakka Kafka still uses the old ActorSystem (Classic System) which is untyped,
      // we need to convert the typed actor system to the classic untyped system.
      implicit val system: ActorSystem = ctx.system.toClassic
      implicit val ec: ExecutionContextExecutor = ctx.executionContext

      AkkaManagement(system).start()
      val consumerProps = ConsumerProperties(system)
      val subscription = Subscriptions.topics(consumerProps.topic)
      /*
      committableSource => offset can be committed after the computation -> at-least-once delivery semantic
      sourceWithOffsetContext => same as committableSource, but with offset position as flow context
      for at-most-once, use atMostOnceSource() and for exactly-once, use plainSource() for a fine-grained control over offset committing
      use Consumer.DrainingControl if you want to stop the stream after some amount of time
       */
      val consume = Consumer.sourceWithOffsetContext(consumerProps.toConsumerSettings, subscription)
        .mapAsync(consumerProps.parallelism) { record => // mapAsync -> process received messages in parallel (batch)
          Future {
            system.log.info(s"consumed record with Id ${record.key()} from partition ${record.partition()}")
            MyMessage.parseFrom(record.value()) // compile with SBT to create MyMessage inside the target folder
            Thread.sleep(consumerProps.sleepTime) // some heavy computation!
          }
        }.runWith(Committer.sinkWithOffsetContext(CommitterSettings(system)))
      // send a message to self if the consumer stops
      consume.onComplete { result =>
        ctx.self ! KafkaConsumerStopped(result)
      }
      Behaviors.receiveMessage[Command] {
        case KafkaConsumerStopped(reason) => //stop the Behavior when the consumer stops
          ctx.log.warn("Consumer stopped {}", reason)
          Behaviors.stopped
      }
    }.narrow //we don't care about the type, so narrow it down to Nothing (instead of Command)
  }
}
