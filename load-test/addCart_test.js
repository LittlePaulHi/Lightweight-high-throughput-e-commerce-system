import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter } from 'k6/metrics';

export const errors = new Counter("errors");
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

  let productid = getRandomInt(10000);
  let quantity = getRandomInt(2000);

  const payload = JSON.stringify({ 'accountID': __VU, 'productID': productid, 'quantity': quantity });
  const params = { headers: { 'Content-Type': 'application/json' }};
  let res = http.post(`${BASE_URL}/api/cart/addCart`, payload, params);

  if(res.status != 200)
    console.log(`[${__VU}] Response status: ${res.status}`);

  const checkRes = check(res, {
    'status is 200': (r) => r.status === 200,
  });

  errors.add(!checkRes);
  
  sleep(500);
}
