import http from 'k6/http';
import { check, group, sleep, fail} from 'k6';
import { Counter } from 'k6/metrics';

export let options = {
  setupTimeout: '10m',
  scenarios: {
      Stage1_getallproducts: {
        executor: 'constant-arrival-rate',
        exec: 'getallproducts',
        rate: 30*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '5m',
        preAllocatedVUs: 25*__ENV.TIMES/10,
        maxVUs: 35*__ENV.TIMES/10
      },
      Stage1_getAllOrderByAccountID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllOrderByAccountID',
        rate: 15*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '2.5m',
        preAllocatedVUs: 10*__ENV.TIMES/10,
        maxVUs: 20*__ENV.TIMES/10
      },
      Stage1_getAllOrderItemsByOrderID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllOrderItemsByOrderID',
        rate: 15*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '3m',
        preAllocatedVUs: 10*__ENV.TIMES/10,
        maxVUs: 20*__ENV.TIMES/10
      },
      Stage1_getAllCartsByAccountID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllCartsByAccountID',
        rate: 15*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '2m',
        preAllocatedVUs: 10*__ENV.TIMES/10,
        maxVUs: 20*__ENV.TIMES/10
      },
      Stage1_addCart: {
        executor: 'constant-arrival-rate',
        exec: 'addCart',
        rate: 11*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '1m',
        preAllocatedVUs: 5*__ENV.TIMES/10,
        maxVUs: 15*__ENV.TIMES/10
      },
      Stage1_editCart: {
        executor: 'constant-arrival-rate',
        exec: 'editCart',
        rate: 11*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '1.5m',
        preAllocatedVUs: 5*__ENV.TIMES/10,
        maxVUs: 15*__ENV.TIMES/10
      },
      Stage1_PurchaseFromCarts: {
        executor: 'constant-arrival-rate',
        exec: 'PurchaseFromCarts',
        rate: 3*__ENV.TIMES/10,
        timeUnit: '1s',
        duration: '2m',
        startTime: '3m',
        preAllocatedVUs: 3*__ENV.TIMES/10,
        maxVUs: 5*__ENV.TIMES/10
      }
  }
};


const BASE_URL = 'http://pp-final.garyxiao.me:3080';
export const errors = new Counter("errors");
let order_data = {};
let cart_data = {};


function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}


export function setup() {
  for (let user = 1; user <= __ENV.TIMES * 5 ; user ++) {
    const params_order = { headers: { 'Content-Type': 'application/json', 'accountID': __VU } };	   
    let res_order = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params_order);

    if(res_order.status == 200) {
      order_data[__VU] = res_order["orders"];
    }

    const params_cart = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res_cart = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_cart);

    if (res_cart.status == 200) {
      cart_data[__VU] = res_cart["cart"];
    }
  }
}


export function getallproducts() {
    
    const res = http.get(`${BASE_URL}/api/product/getAll`);
  
    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });
    
    sleep(500);
}


export function getAllOrderByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU } };	   
    let res = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params);

    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });

    if(res.status == 200) {
      order_data[__VU] = res["orders"];
    }
    
    sleep(500);
}


export function getAllOrderItemsByOrderID() {

    let order = order_data[__VU];
    let orderid = order[getRandomInt(order.length)]["ID"];
  
    const params_getitem = { headers: { 'Content-Type': 'application/json', 'orderID': orderid } };
    let res_getitem = http.get(`${BASE_URL}/api/order/getAllItemsByOrderID`, params_getitem);
  
    const checkRes = check(res_getitem, {
      'status is 200': (r) => r.status === 200,
    });
  
    sleep(500);
}


export function getAllCartsByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params);
  
    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });

    sleep(500);
}


export function addCart() {
    
    let productid = getRandomInt(10000);
    let quantity = getRandomInt(2000);
  
    const payload = JSON.stringify({ 'accountID': __VU, 'productID': productid, 'quantity': quantity });
    const params = { headers: { 'Content-Type': 'application/json' }};
    let res = http.post(`${BASE_URL}/api/cart/addCart`, payload, params);
  
    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });

    if(res.status == 200) {
      let data = JSON.parse(res.body).data;
      
      cart_data[__VU].push(data["cart"]);    
    }  
  
    sleep(500);
}


export function editCart() {

    let cart = cart_data[__VU];
    console.log(cart);
    let cartid = getRandomInt(cart.length);
    let quantity = getRandomInt(2000);
  
    const payload_post = JSON.stringify({ 'accountID': __VU, 'productID': cart[cartid]['ProductID'], 'quantity': quantity, 'cartID': cart[cartid]['ID'] });
    const params_post = { headers: { 'Content-Type': 'application/json' }};
    let res_post = http.post(`${BASE_URL}/api/cart/editCart`, payload_post, params_post);
  
    const checkRes = check(res_post, {
      'status is 200': (r) => r.status === 200,
    });  

    if(res.status == 200) {
      let data = JSON.parse(res.body).data;
      for (let items = 0; items < cart_data[__VU].length; items ++) {
        if(cart_data[__VU][items]["ID"] == data["cart"]["ID"])
          cart_data[__VU].push(data["cart"]);
      }
    }  
  
    sleep(500);
}


export function PurchaseFromCarts() {

    let cart = cart_data[__VU];
    let cartids = [];

    cartids.push(cart['ID']);
  
    const payload_post = JSON.stringify({ 'accountID': __VU, 'cartIDs': cartids });
    const params_post = { headers: { 'Content-Type': 'application/json' } };
    let res_post = http.post(`${BASE_URL}/api/purchase/sync`, payload_post, params_post);  
  
    const checkRes = check(res_post, {
      'status is 200': (r) => r.status === 200,
    });
  
    sleep(500);
}
