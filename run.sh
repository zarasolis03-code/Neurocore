#!/bin/bash
echo "🚀 STARTING NEUROCORE DEPLOYMENT..."
git add .
git commit -m "Update V3.1 - Password Protection Active"
echo "📤 PUSHING TO GITHUB..."
git push -f origin main
echo "✅ DONE! Check your app now."
