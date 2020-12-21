import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter } from 'k6/metrics';

// A simple counter for http requests

export const requests = new Counter('http_reqs');

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

const BASE_URL = 'http://pp-final.garyxiao.me:3080';

export const options = {
  setupTimeout: '10m',
  vusMax: 10000,
  duration: '1m',
  vus: __ENV.TIMES,
  iterations: __ENV.TIMES * 60,
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
  let users = 10000;

  if(__ENV.TIMES * 2 < users)
    users = __ENV.TIMES * 2;

  for (let user = 1; user <= users; user ++) {

    let productid = getRandomInt(10000);
    let quantity = getRandomInt(2000);
  
    const payload = JSON.stringify({ 'accountID': user, 'productID': productid, 'quantity': quantity });
    const params = { headers: { 'Content-Type': 'application/json' }};
    let res_get = http.post(`${BASE_URL}/api/cart/addCart`, payload, params);

    if(res_get.status == 200)
    {
      let data = JSON.parse(res_get.body).data;
      carts[user] = [JSON.parse(JSON.stringify(data["cart"]))];
    }
    else
      carts[user] = [];
  }
  
  sleep(5);
  
  return carts;
}


export default function (data) {

  let cart = data[__VU];

  let cartids = [];

  for (let items = 0; items < cart.length; items ++) {
    cartids.push(cart[items]['ID']);
  }

  const payload_post = JSON.stringify({ 'accountID': __VU, 'cartIDs': cartids });
  const params_post = { headers: { 'Content-Type': 'application/json' } };
  let res_post = http.post(`${BASE_URL}/api/purchase/sync`, payload_post, params_post);

  const checkRes = check(res_post, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.5);
}
