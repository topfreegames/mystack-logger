require 'fluent/test'
require 'fluent/plugin/out_mystack'

class MystackOutputTest < Test::Unit::TestCase
  CONFIG = %[tag foo.bar]

  def setup
    Fluent::Test.setup
    @valid_record = { "kubernetes" => { "container_name" => "mystack-controller" } }
    @invalid_record = { }
    @valid_app_record = { "kubernetes" => { "labels" => { "app" => "foo", "heritage" => "mystack" } } }
    @invalid_app_record = { "kubernetes" => { "labels" => { "foo" => "foo" } } }

    @mystack_output = Fluent::MystackOutput.new
  end

  def test_kubernetes_should_return_true_with_valid_key
    output = Fluent::MystackOutput.new
    assert_true(@mystack_output.kubernetes?(@valid_record))
  end

  def test_kubernetes_should_return_false_with_invalid_key
    assert_false(@mystack_output.kubernetes?(@invalid_record))
  end

  def test_from_container_should_return_true_with_valid_container_name
    assert_true(@mystack_output.from_container?(@valid_record, "mystack-controller"))
  end

  def test_from_container_should_return_false_with_invalid_container_name
    assert_false(@mystack_output.from_container?(@valid_record, "mystack-foo"))
  end

  def test_from_container_should_return_false_with_invalid_record
    assert_false(@mystack_output.from_container?(@invalid_record, "mystack-controller"))
  end

  def test_mystack_deployed_app_should_return_true_with_valid_application_message
    assert_true(@mystack_output.mystack_deployed_app?(@valid_app_record))
  end

  def test_mystack_deployed_app_should_return_false_with_valid_application_message
    assert_false(@mystack_output.mystack_deployed_app?(@invalid_app_record))
  end

end
