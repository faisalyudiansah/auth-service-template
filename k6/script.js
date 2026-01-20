import { check } from "k6";
import http from "k6/http";

const BASE_URL = __ENV.BASE_URL || "http://localhost:8000";

export const options = {
  stages: [
    { duration: "5m", target: 50 }, // ramp-up
    { duration: "5m", target: 100 }, // stable
    { duration: "5m", target: 300 }, // ramp-down
  ],

  summaryTrendStats: ["min", "med", "avg", "count", "max", "p(99)"],

  thresholds: {
    http_req_duration: ["p(99)<500"], // 99% of requests must complete within 500ms
  },
};

export default function () {
  http.get(`${BASE_URL}/products/home?limit=10&page=1`);
}