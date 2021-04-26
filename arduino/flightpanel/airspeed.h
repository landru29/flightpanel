#ifndef __AIRSPEED__H__
#define __AIRSPEED__H__

#include <MCUFRIEND_kbv.h>
#include "image.h"


class AirSpeed {
  public:
    AirSpeed(MCUFRIEND_kbv* tft);
    void init();
    void refresh();
    String name();
    
  private:
    MCUFRIEND_kbv* tft;
    Image *img;

    void drawIndicator(float angle, uint16_t color);

    int16_t ht;
    int16_t top;
    int16_t line;
    int16_t lines;
    int16_t scroll;

    float previousAngle;
};

#endif
