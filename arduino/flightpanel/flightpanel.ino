#include <Adafruit_GFX.h> // Hardware-specific library
#include <MCUFRIEND_kbv.h>
#include "airspeed.h"

MCUFRIEND_kbv tft;

AirSpeed* device;

#define BLACK   0x0000
#define BLUE    0x001F
#define RED     0xF800
#define GREEN   0x07E0
#define CYAN    0x07FF
#define MAGENTA 0xF81F
#define YELLOW  0xFFE0
#define WHITE   0xFFFF

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
    tft.fillScreen(BLACK);
    tft.setTextColor(WHITE, BLACK);
    tft.setTextSize(2);     // System font is 8 pixels.  ht = 8*2=16
    tft.setCursor(100, 0);
    tft.print("ID = 0x");
    tft.println(id, HEX);
    tft.setCursor(0, 0);
}



void loop()
{
  device->refresh();
}
