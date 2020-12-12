import http from 'k6/http';
import { check, group, sleep, fail } from 'k6';

export let options = {
  scenarios: {
      Stage1_getallproducts: {
        executor: 'constant-arrival-rate',
        exec: 'getallproducts',
        rate: 30*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//	startTime: '16m',
        preAllocatedVUs: 25*__ENV.TIMES/100,
        maxVUs: 35*__ENV.TIMES/100
      },
      Stage1_getAllOrderByAccountID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllOrderByAccountID',
        rate: 15*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 10*__ENV.TIMES/100,
        maxVUs: 20*__ENV.TIMES/100
      },
      Stage1_getAllOrderItemsByOrderID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllOrderItemsByOrderID',
        rate: 15*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 10*__ENV.TIMES/100,
        maxVUs: 20*__ENV.TIMES/100
      },
      Stage1_getAllCartsByAccountID: {
        executor: 'constant-arrival-rate',
        exec: 'getAllCartsByAccountID',
        rate: 15*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 10*__ENV.TIMES/100,
        maxVUs: 20*__ENV.TIMES/100
      },
      Stage1_addCart: {
        executor: 'constant-arrival-rate',
        exec: 'addCart',
        rate: 11*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 5*__ENV.TIMES/100,
        maxVUs: 15*__ENV.TIMES/100
      },
      Stage1_editCart: {
        executor: 'constant-arrival-rate',
        exec: 'editCart',
        rate: 11*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 5*__ENV.TIMES/100,
        maxVUs: 15*__ENV.TIMES/100
      },
      Stage1_PurchaseFromCarts: {
        executor: 'constant-arrival-rate',
        exec: 'PurchaseFromCarts',
        rate: 3*__ENV.TIMES/100,
        timeUnit: '1s',
        duration: '1m',
//        startTime: '16m',
        preAllocatedVUs: 3*__ENV.TIMES/100,
        maxVUs: 5*__ENV.TIMES/100
      }
  }
};


const BASE_URL = 'http://pp-final.garyxiao.me:3080';


function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}


export function getallproducts() {
    
    let res = http.get(`${BASE_URL}/api/product/getAll`);
    check(res, { 'status was 200': (r) => r.status == 200 });

    const checkRes = check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(500);
}


export function getAllOrderByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU } };
    let res = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params);

    if(res.status != 200)
         console.log(`[${__VU}] Response status: ${res.status}`);    

    const checkRes = check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(500);
}


export function getAllOrderItemsByOrderID() {

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


export function getAllCartsByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params);

    if(res.status != 200)
        console.log(`[${__VU}] Response status: ${res.status}`);
    
    const checkRes = check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(500);
}


export function addCart() {

    let productid = getRandomInt(20000);
    const payload = JSON.stringify({ 'accountID': __VU, 'productID': productid, 'quantity': 1 });
    const params = { headers: { 'Content-Type': 'application/json' }};
    let res = http.post(`${BASE_URL}/api/cart/addCart`, payload, params);

    if(res.status != 200)
        console.log(`[${__VU}] Response status: ${res.status}`);
   
    const checkRes = check(res, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(500);
}


export function editCart() {

    const params_get = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res_get = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_get);
    let data = JSON.parse(res_get.body).data;
    let cart = data["cart"];
    let cartid;
    let quantity;

    if(cart.length == 0) {
        check(res_get, { 'status is 200': (r) => r.status === 200, });
        return;
    }
    else {
        cartid = getRandomInt(cart.length);
        quantity = getRandomInt(2000);
    }

    sleep(100);
}


export function PurchaseFromCarts() {

    const params_get = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res_get = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_get);
    let data = JSON.parse(res_get.body).data;
    let cart = data["cart"];
    let cartids = [];

    if(cart.length == 0) {
        check(res_get, { 'status is 200': (r) => r.status === 200, });
        return;
    }
    else {
        let num_of_carts = getRandomInt(cart.length);
        for (let step = 0; step < num_of_carts; step ++) {
            cartids.push(cart[getRandomInt(cart.length)]['ID']);
        }
    }

    sleep(100);

    const payload_post = JSON.stringify({ 'accountID': __VU, 'cartIDs': cartids });
    const params_post = { headers: { 'Content-Type': 'application/json' } };
    let res_post = http.post(`${BASE_URL}/api/purchase/sync`, payload_post, params_post);

    if(res_post.status != 200)
        console.log(`[${__VU}] Response status: ${res_post.status}`);

    const checkRes = check(res_post, {
        'status is 200': (r) => r.status === 200,
    });

    sleep(300);
}
