#include "airspeed.h"
#include "communication.h"
//#include <serial>


AirSpeed::AirSpeed(MCUFRIEND_kbv* tft) {
  this->tft = tft;

  this->ht = 16;
  this->top = 1;
  this->lines = 20;

  for (this->line = 1; this->line < 21; this->line++) {
    this->tft->println(String(this->line) + ": ");
  }
}

String AirSpeed::name() {
  return String("GYRO");
}

void AirSpeed::refresh() {
  float value = readFloat();
  this->tft->setCursor(0, (this->scroll + this->top) * this->ht);
  if (++this->scroll >= this->lines) {
    this->scroll = 0;
  }
  this->tft->vertScroll(this->top * this->ht, this->lines * this->ht, this->scroll * this->ht);
  this->tft->println(String(this->line) + ": " + value);
  delay(100);
  this->line++;
}
