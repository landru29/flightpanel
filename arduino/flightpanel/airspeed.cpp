#include "airspeed.h"
#include "communication.h"
#include <Arduino.h>


AirSpeed::AirSpeed(MCUFRIEND_kbv* tft) {
  this->tft = tft;
  this->img = new Image(tft);

  this->ht = 16;
  this->top = 1;
  this->lines = 20;

  for (this->line = 1; this->line < 21; this->line++) {
    this->tft->println(String(this->line) + ": ");
  }

  this->previousAngle = 0.0;
}

String AirSpeed::name() {
  return String("AIR");
}

void AirSpeed::refresh() {
  float value = readFloat();
  this->drawIndicator(this->previousAngle, TFT_BLACK);
  this->drawIndicator(value, TFT_WHITE);
  this->previousAngle = value;
  /*this->tft->setCursor(0, (this->scroll + this->top) * this->ht);
  if (++this->scroll >= this->lines) {
    this->scroll = 0;
  }
  this->tft->vertScroll(this->top * this->ht, this->lines * this->ht, this->scroll * this->ht);
  this->tft->println(String(this->line) + ": " + value);
  delay(100);
  this->line++;*/
  
}

void AirSpeed::init() {
  this->img->showBMP("/airspeed.bmp", 5, 5);
}

void AirSpeed::drawIndicator(float speed, uint16_t color) {
  int xOffset=159, yOffset=159;
  float angle = speed * 1.75-30;
  this->tft->drawLine(xOffset, yOffset, xOffset+80*sin(3.14*angle/180), yOffset-80*cos(3.14*angle/180), color);
}
