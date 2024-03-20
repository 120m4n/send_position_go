@echo off
send_position.exe -json=trazado.json -url=http://localhost:53430/api/coordinates/ -parametro=918422f729de0348 -puntos=13 -ciclico=true -pausa=10000 -verbose=true