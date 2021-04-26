#include <Adafruit_GFX.h> // Hardware-specific library
#include <MCUFRIEND_kbv.h>
#include "airspeed.h"

#if defined(ESP32)
#define SD_CS     5
#else
#define SD_CS     10
#endif

MCUFRIEND_kbv tft;

AirSpeed* device;


void setup()
{
    tft.reset();
    uint16_t id = tft.readID();

    device = new AirSpeed(&tft);
    
    Serial.begin(115200);

    String str = device->name();
    str += "/" + String(id, HEX);

    Serial.println(str);
  
    
    tft.begin(id);
    tft.setRotation(0);     //Portrait
    tft.fillScreen(TFT_BLACK);
    tft.setTextColor(TFT_WHITE, TFT_BLACK);
    tft.setTextSize(2);     // System font is 8 pixels.  ht = 8*2=16
    tft.setCursor(100, 0);
    //tft.print("ID = 0x");
    //tft.println(id, HEX);
    tft.setCursor(0, 0);

     bool good = SD.begin(SD_CS);
    if (!good) {
        Serial.print(F("cannot start SD"));
        while (1);
    }

    device->init();
}



void loop()
{
  device->refresh();
}
