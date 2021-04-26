#ifndef __AIRSPEED__H__
#define __AIRSPEED__H__

#include <MCUFRIEND_kbv.h>


class AirSpeed {
  public:
    AirSpeed(MCUFRIEND_kbv* tft);
    void refresh();
    String name();
    
  private:
    MCUFRIEND_kbv* tft;

    int16_t ht;
    int16_t top;
    int16_t line;
    int16_t lines;
    int16_t scroll;
};

#endif
