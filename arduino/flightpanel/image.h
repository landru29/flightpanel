#ifndef __IMAGE_H__
#define __IMAGE_H__

#define BUFFPIXEL      20
#define PALETTEDEPTH   0     // do not support Palette modes

#include <MCUFRIEND_kbv.h>
#include <Arduino.h>
#include <SD.h>  


class Image {
  public:
    Image(MCUFRIEND_kbv* tft);
    uint8_t showBMP(char *nm, int x, int y);
    
  private:
    MCUFRIEND_kbv* tft;

    uint16_t read16(File& f);
    uint32_t read32(File& f);

};



#endif
