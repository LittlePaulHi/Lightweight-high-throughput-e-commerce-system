import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter } from 'k6/metrics';

export const requests = new Counter('http_reqs');

const BASE_URL = 'http://pp-final.garyxiao.me:3080';

export const options = {
  stages: [
    { target: __ENV.TIMES, duration: '30s' },
    { target: __ENV.TIMES, duration: '1m' },
    { target: __ENV.TIMES, duration: '30s' },
  ],
  thresholds: {
    Errors: ['count < 10'],
    http_req_duration: ['avg < 2000'],
    http_req_waiting: ['avg < 1000'],
  },
};

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}

export default function () {
  
  const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
  let res = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params);

  const checkRes = check(res, {
    'status is 200': (r) => r.status === 200,
  });
  
  sleep(500);
}
