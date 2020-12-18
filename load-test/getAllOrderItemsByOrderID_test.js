import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter } from 'k6/metrics';

export const requests = new Counter('http_reqs');

const BASE_URL = 'http://pp-final.garyxiao.me:3080';

export const options = {
  setupTimeout: '10m',
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

  let orders = {};

  for(let user=1; user <= __ENV.TIMES; user++) {
    
    const params_get = { headers: { 'Content-Type': 'application/json', 'accountID': user } };
    let res_get = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params_get);  
  
    if (res_get.status == 200) {
      let data = JSON.parse(res_get.body).data;
      orders[user] = data["orders"];
    }
  }

  return orders;
}

export default function (data) {

  let order = data[__VU];
  
  let orderid = order[getRandomInt(0)]["ID"];

  const params_getitem = { headers: { 'Content-Type': 'application/json', 'orderID': orderid } };
  let res_getitem = http.get(`${BASE_URL}/api/order/getAllItemsByOrderID`, params_getitem);

  const checkRes = check(res_getitem, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(300);
}
