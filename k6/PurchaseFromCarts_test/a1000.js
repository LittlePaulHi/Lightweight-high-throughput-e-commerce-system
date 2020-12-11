import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter } from 'k6/metrics';

// A simple counter for http requests

export const requests = new Counter('http_reqs');

// you can specify stages of your test (ramp up/down patterns) through the options object
// target is the number of VUs you are aiming for

const BASE_URL = 'http://pp-final.garyxiao.me:3080';

export const options = {
  stages: [
    { target: 1000, duration: '30s' },
    { target: 1000, duration: '1m' },
    { target: 1000, duration: '30s' },
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
  // our HTTP request, note that we are saving the response to res, which can be accessed later
  const params_get = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
  let res_get = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_get);

  let data = JSON.parse(res_get.body).data;

  let cart = data["cart"];
  
  let cartids = [];

  if(cart.length == 0) {
    check(res_get, { 'status is 200': (r) => r.status === 200, });
    console.log('Length?' );
    return;
  }
  else {
    let num_of_carts = getRandomInt(cart.length);
    for (let step = 0; step < num_of_carts; step ++) {
      cartids.push(cart[getRandomInt(cart.length)]['ID']);
    }
  }

  sleep(100);

  console.log('Change ' + __VU + ' CartIDs: ' + cartids);

  const payload_post = JSON.stringify({ 'accountID': __VU, 'cartIDs': cartids });
  const params_post = { headers: { 'Content-Type': 'application/json' } };
  let res_post = http.post(`${BASE_URL}/api/purchase/sync`, payload_post, params_post);  

  console.log('success!');

  if(res_post.status != 200)
    console.log(`[${__VU}] Response status: ${res_post.status}`);

  const checkRes = check(res_post, {
    'status is 200': (r) => r.status === 200,
  });
  
  sleep(300);
}
