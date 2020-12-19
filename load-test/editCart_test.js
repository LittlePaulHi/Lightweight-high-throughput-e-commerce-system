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

export function setup() {

  let carts = {};

  for (let user = 1; user <= __VU.TIMES; user ++) {
    
    let params_get = { headers: { 'Content-Type': 'application/json', 'accountID': user, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res_get = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_get);
    let data = JSON.parse(res_get.body).data;
  
    carts[user] = data["cart"];
  }
  return carts;
}

export default function (data) {

  let cart = data[__VU];

  let cartid = getRandomInt(cart.length);
  let quantity = getRandomInt(2000);

  const payload_post = JSON.stringify({ 'accountID': __VU, 'productID': cart[cartid]['ProductID'], 'quantity': quantity, 'cartID': cart[cartid]['ID'] });
  const params_post = { headers: { 'Content-Type': 'application/json' }};
  let res_post = http.post(`${BASE_URL}/api/cart/editCart`, payload_post, params_post);

  const checkRes = check(res_post, {
    'status is 200': (r) => r.status === 200,
  });  

  sleep(500);
}
