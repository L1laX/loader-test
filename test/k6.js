import http from "k6/http"
import { sleep } from "k6"

export let options = {
  // Simulate 100,000 RPM (1,667 RPS)
  vus: 2000, // Number of virtual users
  duration: "1m", // Duration of the test (1 minute)
}

export default function () {
  http.get("http://localhost/check/12345") // Replace with your API endpoint
  sleep(0) // No delay between requests
}
