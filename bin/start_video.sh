cd ~/Development/mjpg-streamer/mjpg-streamer-experimental/_build
export LD_LIBRARY_PATH=/usr/local/lib/mjpg-streamer
sudo ./mjpg_streamer -i "/usr/local/lib/mjpg-streamer/input_uvc.so -d /dev/video1" -o "/usr/local/lib/mjpg-streamer/output_http.so -w ../www"
cd ~/Preso/endesa2016
