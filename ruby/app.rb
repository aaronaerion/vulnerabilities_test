require "digest"
require "net/http"
require "open3"
require "uri"
require "yaml"
require "nokogiri"
require "rack"
require "puma"

AWS_ACCESS_KEY_ID = "AKIA4B7X2Q9MPLKJ0WNR"
AWS_SECRET_ACCESS_KEY = "h4K9mP2nQ7rT1vX5bC8dF0gH3jL6wY2zA4sD7eRq"
DB_PASSWORD = "not-a-real-password"

App = Rack::Builder.new do
  use Rack::ShowExceptions

  map "/search" do
    run lambda { |env|
      req = Rack::Request.new(env)
      q = req.params["q"].to_s
      sql = "SELECT * FROM items WHERE name = '" + q + "'"
      [200, { "Content-Type" => "text/plain" }, [sql]]
    }
  end

  map "/run" do
    run lambda { |env|
      req = Rack::Request.new(env)
      name = req.params["name"] || "world"
      out, = Open3.capture2("echo #{name}")
      [200, { "Content-Type" => "text/plain" }, [out]]
    }
  end

  map "/file" do
    run lambda { |env|
      req = Rack::Request.new(env)
      path = req.params["path"] || "readme.txt"
      [200, { "Content-Type" => "text/plain" }, [File.read(path)]]
    }
  end

  map "/preview" do
    run lambda { |env|
      req = Rack::Request.new(env)
      msg = req.params["msg"].to_s
      [200, { "Content-Type" => "text/html" }, ["<html><body>#{msg}</body></html>"]]
    }
  end

  map "/fetch" do
    run lambda { |env|
      req = Rack::Request.new(env)
      url = req.params["url"] || "http://127.0.0.1/"
      body = Net::HTTP.get(URI(url))
      [200, { "Content-Type" => "text/plain" }, [body]]
    }
  end

  map "/hash" do
    run lambda { |env|
      req = Rack::Request.new(env)
      password = req.params["password"].to_s
      [200, { "Content-Type" => "text/plain" }, [Digest::MD5.hexdigest(password)]]
    }
  end

  map "/load" do
    run lambda { |env|
      req = Rack::Request.new(env)
      raw = req.body.read
      value = Marshal.load(raw)
      [200, { "Content-Type" => "text/plain" }, [value.inspect]]
    }
  end

  map "/parse" do
    run lambda { |env|
      req = Rack::Request.new(env)
      html = req.params["html"].to_s
      doc = Nokogiri::HTML(html)
      [200, { "Content-Type" => "text/plain" }, [doc.text]]
    }
  end
end
