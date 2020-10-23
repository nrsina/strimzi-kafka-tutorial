package com.snourian.strimzi.consumer

import akka.actor.typed.ActorSystem

object Application extends App {
  ActorSystem[Nothing](ConsumerClient(), "StrimziConsumer")
}
