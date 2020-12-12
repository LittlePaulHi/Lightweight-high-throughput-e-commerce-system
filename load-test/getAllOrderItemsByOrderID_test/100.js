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
    { target: 100, duration: '30s' },
    { target: 100, duration: '1m' },
    { target: 100, duration: '30s' },
  ],
  thresholds: {
    Errors: ['count < 10'],
    http_req_duration: ['avg < 2000'],
    http_req_waiting: ['avg < 1000'],
  },
};

export default function () {
  // our HTTP request, note that we are saving the response to res, which can be accessed later
  const params_get = { headers: { 'Content-Type': 'application/json', 'accountID': __VU } };
  let res_get = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params_get);  

  let data = JSON.parse(res_get.body).data;

  let order = data["orders"];

  let orderid;
  let quantity;

  if(order.length == 0) {
    check(res_get, { 'status is 200': (r) => r.status === 200, });
    sleep(500);
    console.log('No data');
    return;
  }
  else {
    orderid = getRandomInt(order.length);
    quantity = getRandomInt(2000);
  }

  sleep(100);

  const params_getitem = { headers: { 'Content-Type': 'application/json', 'orderID': orderid } };
  let res_getitem = http.get(`${BASE_URL}/api/order/getAllItemsByOrderID`, params_getitem);

  if(res_getitem.status != 200)
    console.log(`[${__VU}] Response status: ${res_getitem.status}`);

  const checkRes = check(res_getitem, {
    'status is 200': (r) => r.status === 200,
  });
  
  sleep(500);
}
