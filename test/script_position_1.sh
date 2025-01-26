#!/bin/bash

./send_position -url=https://traceo.test-electrosoftware.xyz/coordinates -puntos=33 -fleet=avatar -userid=Gtest/test_1 -uniqueid=test1 -verbose=true -lat=3.431239 -lon=-76.541954118090 &
./send_position -url=https://traceo.test-electrosoftware.xyz/coordinates -puntos=33  -fleet=avatar -userid=Gtest/test_2 -uniqueid=test2 -verbose=true -lat=3.431239 -lon=-76.541954118090 &
./send_position -url=https://traceo.test-electrosoftware.xyz/coordinates -puntos=33 -fleet=avatar -userid=Gtest/test_3 -uniqueid=test3 -verbose=true -lat=3.431239 -lon=-76.541954118090 &
./send_position -url=https://traceo.test-electrosoftware.xyz/coordinates -puntos=33 -fleet=avatar -userid=Gtest/test_4 -uniqueid=test4 -verbose=true -lat=3.431239 -lon=-76.541954118090 &
./send_position -url=https://traceo.test-electrosoftware.xyz/coordinates -puntos=33 -fleet=avatar -userid=Gtest/test_5 -uniqueid=test5 -verbose=true -lat=3.431239 -lon=-76.541954118090 &

# Wait for all background jobs to finish
wait