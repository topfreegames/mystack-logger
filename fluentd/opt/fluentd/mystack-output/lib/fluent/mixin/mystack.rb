require 'json'
require 'nsq'
require 'influxdb'
require 'yajl/json_gem'

module Fluent
  module Mixin
    module Mystack
      LOGGER_URL = "http://#{ENV['MYSTACK_LOGGER_SERVICE_HOST']}:#{ENV['MYSTACK_LOGGER_SERVICE_PORT_HTTP']}/logs"
      NSQ_URL = "#{ENV['MYSTACK_NSQD_SERVICE_HOST']}:#{ENV['MYSTACK_NSQD_SERVICE_PORT_TRANSPORT']}"

      def kubernetes?(message)
        return message["kubernetes"] != nil
      end

      def from_controller?(message)
        if from_container?(message, "mystack-controller")
          return message["log"] =~ /^(INFO|WARN|DEBUG|ERROR)\s+(\[(\S+)\])+:(.*)/
        end
        return false
      end

      def from_container?(message, regex)
        if kubernetes? message
          return true if Regexp.new(regex).match(message["kubernetes"]["container_name"]) != nil
        end
        return false
      end

      def mystack_deployed_app?(message)
        if kubernetes? message
          labels = message["kubernetes"]["labels"]
          return true if message["kubernetes"]["namespace_name"] != "mystack" && labels["heritage"] == "mystack" && labels["app"] != nil
        end
        return false
      end

      def push(producer, value)
        begin
          if value.kind_of? Hash
            producer.write(JSON.dump(value))
          else
            producer.write(value)
          end
        rescue Exception => e
          puts "Error:#{e.message}"
          puts e.backtrace
        end
      end

      def get_nsq_producer(topic)
        begin
          puts "Creating nsq producer (#{NSQ_URL}) for topic:#{topic}"
          return Nsq::Producer.new(nsqd: NSQ_URL, topic: topic)
        rescue Exception => e
          puts "Error:#{e.message}"
          return nil
        end
      end
    end
  end
end
