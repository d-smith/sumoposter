## Sumo Poster

Post messages to a hosted Sumologic collector via an integration with Fluentd

### Overview

My application emits log messages that are formatted using JSON. If I use
the [fluentd Sumologic plugin](https://github.com/memorycraft/fluent-plugin-sumologic), the plugin bufffers up messages, then sends a single giant message to sumo escaping all buffered data into a single string with new lines separating message components. Here's an example:

<pre>
opening connection to xxx.xxx.xxx.xxx:8000... 
opened 
starting SSL for xxx.xxx.xxx.xxx:8000... 
SSL established 
<- "CONNECT endpoint1.collection.us2.sumologic.com:443 HTTP/1.1\r\nHost: endpoint1.collection.us2.sumologic.com:443\r\n\r\n" 
-> "HTTP/1.1 200 Connection established\r\n" 
-> "\r\n" 
<- "POST /receiver/v1/http/my-secret-sauce HTTP/1.1\r\nAccept-Encoding: gzip;q=1.0,deflate;q=0.6,identity;q=0.3\r\nAccept: */*\r\nUser-Agent: Ruby\r\nConnection: close\r\nHost: endpoint1.collection.us2.sumologic.com\r\nContent-Length: 1445\r\nContent-Type: application/x-www-form-urlencoded\r\n\r\n" 
<- "{\"level\":\"info\",\"msg\":\"session: 662272250 symbol xxx\",\"time\":\"2016-02-01T23:19:57Z\"}\n{\"Name\":\"Quote\",\"Duration\":17028440,\"time\":\"2016-02-01T23:19:57.512958787Z\",\"TxnId\":\"6733ecea-1973-e4ed-d4d9-1f53b97c2a5a\",\"Contributors\":[{\"Name\":\"quote svc plugin\",\"Duration\":16973253,\"Error\":\"\",\"ServiceCalls\":null},{\"Name\":\"quote-backend backend\",\"Duration\":4147820,\"Error\":\"\",\"ServiceCalls\":[{\"Name\":\"backend call localhost:4545\",\"Duration\":4062620,\"Error\":\"\"}]}],\"ErrorFree\":true,\"Error\":\"\"}\n{\"level\":\"info\",\"msg\":\"session: 437994385 symbol xxx\",\"time\":\"2016-02-01T23:19:58Z\"}\n{\"Name\":\"Quote\",\"Duration\":15330902,\"time\":\"2016-02-01T23:19:58.551407994Z\",\"TxnId\":\"a547094e-d088-4a3e-4ab9-d5a99fa49f94\",\"Contributors\":[{\"Name\":\"quote svc plugin\",\"Duration\":15321536,\"Error\":\"\",\"ServiceCalls\":null},{\"Name\":\"quote-backend backend\",\"Duration\":3657521,\"Error\":\"\",\"ServiceCalls\":[{\"Name\":\"backend call localhost:4545\",\"Duration\":3603709,\"Error\":\"\"}]}],\"ErrorFree\":true,\"Error\":\"\"}\n{\"level\":\"info\",\"msg\":\"session: 851848122 symbol xxx\",\"time\":\"2016-02-01T23:19:59Z\"}\n{\"Name\":\"Quote\",\"Duration\":91256478,\"time\":\"2016-02-01T23:19:59.092173739Z\",\"TxnId\":\"0f85d9b1-b606-7ce7-493f-a2ab2f896231\",\"Contributors\":[{\"Name\":\"quote svc plugin\",\"Duration\":91248048,\"Error\":\"\",\"ServiceCalls\":null},{\"Name\":\"quote-backend backend\",\"Duration\":3513933,\"Error\":\"\",\"ServiceCalls\":[{\"Name\":\"backend call localhost:4545\",\"Duration\":3389380,\"Error\":\"\"}]}],\"ErrorFree\":true,\"Error\":\"\"}" 
-> "HTTP/1.1 200 OK\r\n" 
-> "Cache-control: no-cache=\"set-cookie\"\r\n" 
-> "Date: Mon, 01 Feb 2016 23:20:54 GMT\r\n" 
-> "Set-Cookie: AWSELB=93EB03E50CD3C318B325E174BF375A879C9EB23656CC5DB42F4090D7A3077C838B1734CF3DF0AC46F2407B7CAB091D2E69B2D2A919ECF47B52A194D678630A20D8D6443031;PATH=/\r\n"
-> "Strict-Transport-Security: max-age=15552000\r\n" 
-> "X-Content-Type-Options: nosniff\r\n" 
-> "X-Frame-Options: SAMEORIGIN\r\n" 
-> "X-XSS-Protection: 1; mode=block\r\n" 
-> "Content-Length: 0\r\n" 
-> "Connection: Close\r\n" 
-> "\r\n" 
reading 0 bytes... 
-> "" 
read 0 bytes 
Conn close
</pre>

To get around this, I wrote my own message poster, which can be integrated into fluend like this:

```xml
<source>
  type tail
  path /home/vagrant/goxavi/src/github.com/xtracdev/xavisample/xs.log
  pos_file ./x-log-pos
  tag xavi
  time_key fake
  format json
</source>


<match **>
  @type copy

  <store>
  @type exec
  command ./sumoposter
  buffer_path ./bp
  format json
  flush_interval 5s
  </store>

  <store>
  @type stdout
  </store>
</match>
```


Note in the above configuration the stdout output is used just for debugging purposes - it would be removed for normal usage.

