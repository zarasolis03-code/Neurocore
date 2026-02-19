#!/data/data/com.termux/files/usr/bin/bash
echo "Installing Neurocore App..."
curl -L -o neuro_core https://github.com/zarasolis03-code/Neurocore/raw/main/neuro_core
chmod +x neuro_core
echo "Done! Type ./neuro_core to start."
./neuro_core
