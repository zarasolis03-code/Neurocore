let basePrice = 0.000512405;
let totalMined = 20800; // Твојот моментален баланс

function updateDashboard() {
    // Симулираме раст врз основа на нови блокови
    totalMined += Math.floor(Math.random() * 50);
    
    // Алгоритам: Како што расте ископаната сума, цената малку варира (како на берза)
    let volatility = (Math.random() * 0.02) - 0.01; // +/- 1% промена
    basePrice = basePrice + (basePrice * volatility);

    document.getElementById('price').innerText = "$" + basePrice.toFixed(6);
    document.getElementById('blocks').innerHTML += "<p>✅ Block sync: Network Difficulty 5 | Mined: " + totalMined + " NEURO</p>";
    
    // Автоматски скрол најдолу во фидот
    const blocksDiv = document.getElementById('blocks');
    blocksDiv.scrollTop = blocksDiv.scrollHeight;
}

// Освежувај на секои 3 секунди
setInterval(updateDashboard, 3000);
