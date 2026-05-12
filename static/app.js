async function loadData() {
    const response = await fetch("/reading")
    const data = await response.json()

    const labels = data.map(r => r.timestamp)
    const values = data.map(r => r.co2)

    const ctx = document.getElementById("chart")

    new Chart(ctx, {
        type: "line",
        data: {
            labels: labels,
            datasets: [{
                label: "CO2",
                data: values
            }]
        }
    })
}

loadData()