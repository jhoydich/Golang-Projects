#include <PubSubClient.h>
#include <ESP8266WiFi.h>

const int LED = 2;
String ssid = "";
String pwd = "";
const char* mqtt_server = "broker.hivemq.com";

WiFiClient espClient;
PubSubClient client(espClient);

void wifiConnect(String ssid, String pwd);
void callback(char* topic, byte* payload, unsigned int length);
void reconnect();

void setup() {
  // put your setup code here, to run once:
  Serial.begin(9600);
  wifiConnect(ssid, pwd);
  delay(2000);
  //tcp connexction to mqttserver
  client.setServer(mqtt_server, 1883);

  client.setCallback(callback);

  while(!client.connected()) {
    String clientid = "JSHESP8266";
    Serial.println("connecting to hivemq");
    if (client.connect(clientid.c_str(),"username","password","esp/led",0,true,"last will")) {
      Serial.println("Connected");
    } else {
      Serial.print("failed =");
      Serial.println(client.state());
      delay(5000);
    }
  }
}

unsigned long t = millis();

void loop() {
  // put your main code here, to run repeatedly:
  unsigned long t2 = millis();
  if (!client.connected()){
    reconnect();
  }
  client.loop();
  if (t2 - t > 2000){
    client.publish("jsh/roomtemp","60");
    client.publish("jsh/roommoisture", "75");
    client.publish("jsh/soiltemp", "54");
    client.publish("jsh/soilmoisture", "45.4");
    
    t = millis();
  }
  
  delay(1000);

}

void wifiConnect(String ssid, String pwd){
 
  WiFi.begin(ssid, pwd);   // add Wi-Fi networks you want to connect to
 
  Serial.println("Connecting ...");
  
  while (WiFi.status() != WL_CONNECTED) { // Wait for the Wi-Fi to connect: scan for Wi-Fi networks, and connect to the strongest of the networks above
    delay(250);
    Serial.print('.');
  }
  
  Serial.println('\n');
  Serial.print("Connected to ");
  Serial.println(WiFi.SSID());              // Tell us what network we're connected to
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());           // Send the IP address of the ESP8266 to the computer 
}

void callback(char* topic, byte* payload, unsigned int length){

  if (strcmp(topic,"esp/led") == 0 &&(char) payload[0]== '1'){
      digitalWrite(LED, HIGH);
      //Serial.println("ON");
  } else {
      digitalWrite(LED, LOW);
      //Serial.println("OFF");
  }
  

  
}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    // Create a random client ID
    String clientId = "JSHESP";
   
    // Attempt to connect
    if (client.connect(clientId.c_str())) {
      Serial.println("connected");
 
      client.subscribe("esp/led");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}
