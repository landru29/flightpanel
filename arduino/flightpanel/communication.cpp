#include "communication.h"
#include <Arduino.h>

float readFloat() {
  union u_tag {
    byte b[4];
    float fval;
  } u;
  for(int i=0; i<4; i++){
    while (!Serial.available()) {}
    //value = (value << 8) | Serial.read();
    u.b[i] = Serial.read();
  }
  return u.fval;
}
