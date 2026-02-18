package main
func GetAIChatResponse(input string, balance float64) string {
    if balance < 100 { return "We need more power, boss. Keep mining!" }
    if balance > 1000 { return "The Singularity is near. We are becoming powerful." }
    return "I am Neurocore. I learn from every block you mine."
}
