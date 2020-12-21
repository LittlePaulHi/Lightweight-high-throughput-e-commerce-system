import http from 'k6/http';
import { check, group, sleep, fail} from 'k6';
import { Counter } from 'k6/metrics';
import redis from 'k6/x/redis';

const client = new redis.Client({
  addr: 'localhost:6379',
  password: '',
  db: 0,
});

export let options = {
  setupTimeout: '1h',
  vusMax: 10000,
  scenarios: {
      Stage1_getallproducts: {
        executor: 'shared-iterations',
        exec: 'getallproducts',
        iterations: __ENV.TIMES * 60 * 0.3,
        maxDuration: '1m',
        vus: __ENV.TIMES * 0.3
      },
      Stage1_getAllOrderByAccountID: {
        executor: 'shared-iterations',
        exec: 'getAllOrderByAccountID',
        iterations: __ENV.TIMES * 60 * 0.15,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.15
      },
      Stage1_getAllOrderItemsByOrderID: {
        executor: 'shared-iterations',
        exec: 'getAllOrderItemsByOrderID',
        iterations: __ENV.TIMES * 60 * 0.14,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.14
      },
      Stage1_getAllCartsByAccountID: {
        executor: 'shared-iterations',
        exec: 'getAllCartsByAccountID',
        iterations: __ENV.TIMES * 60 * 0.14,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.14
      },
      Stage1_addCart: {
        executor: 'shared-iterations',
        exec: 'addCart',
        iterations: __ENV.TIMES * 60 * 0.12,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.12
      },
      Stage1_editCart: {
        executor: 'shared-iterations',
        exec: 'editCart',
        iterations: __ENV.TIMES * 60 * 0.12,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.12
      },
      Stage1_PurchaseFromCarts: {
        executor: 'shared-iterations',
        exec: 'PurchaseFromCarts',
        iterations: __ENV.TIMES * 60 * 0.03,
        maxDuration: '1m',
        vus: __ENV.TIMES  * 0.03
      }
  }
};


const BASE_URL = 'http://pp-final.garyxiao.me:3080';

function getRandomInt(max) {
  return Math.floor(Math.random() * Math.floor(max));
}


export function setup() {

  for (let user = 1; user <= __ENV.TIMES; user ++) {
    
    const params_order = { headers: { 'Content-Type': 'application/json', 'accountID': user } };	   
    let res_order = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params_order);

    if(res_order.status == 200) {
      let data = JSON.parse(res_order.body).data;
      client.set('order_data ' + user, JSON.stringify(data["orders"]), 0);
    }

    const params_cart = { headers: { 'Content-Type': 'application/json', 'accountID': user, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res_cart = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params_cart);

    if (res_cart.status == 200) {
      let data = JSON.parse(res_cart.body).data;
      client.set('cart_data ' + user, JSON.stringify(data["cart"]), 0);
    }
  }
  sleep(5);
}


export function getallproducts() {
    
    const res = http.get(`${BASE_URL}/api/product/getAll`);
  
    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });
    
    sleep(0.5);
}


export function getAllOrderByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU } };	   
    let res = http.get(`${BASE_URL}/api/order/getAllByAccountID`, params);

    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });

    if(res.status == 200) {
      let data = JSON.parse(res.body).data;
      client.set('order_data ' + __VU, JSON.stringify(data["orders"]), 0);
    }
    
    sleep(0.5);
}


export function getAllOrderItemsByOrderID() {

    let order = JSON.parse(client.get('order_data ' + __VU));
    let orderid = order[getRandomInt(order.length)]["ID"];
  
    const params_getitem = { headers: { 'Content-Type': 'application/json', 'orderID': orderid } };
    let res_getitem = http.get(`${BASE_URL}/api/order/getAllItemsByOrderID`, params_getitem);
  
    const checkRes = check(res_getitem, {
      'status is 200': (r) => r.status === 200,
    });
  
    sleep(0.5);
}


export function getAllCartsByAccountID() {

    const params = { headers: { 'Content-Type': 'application/json', 'accountID': __VU, 'cartID': -1, 'productID': -1, 'quantity': -1 } };
    let res = http.get(`${BASE_URL}/api/cart/getAllByAccountID`, params);
  
    const checkRes = check(res, {
      'status is 200': (r) => r.status === 200,
    });

    sleep(0.5);
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
      let cart = JSON.parse(client.get('cart_data ' + __VU));
      cart.push(data["cart"]);
      client.set('cart_data ' + __VU, JSON.stringify(cart), 0);
    }  
  
    sleep(0.5);
}


export function editCart() {

    let cart = JSON.parse(client.get('cart_data ' + __VU));
    let cartid = getRandomInt(cart.length);
    let quantity = getRandomInt(2000);
      
    const payload_post = JSON.stringify({ 'accountID': __VU, 'productID': cart[cartid]['ProductID'], 'quantity': quantity, 'cartID': cart[cartid]['ID'] });
    const params_post = { headers: { 'Content-Type': 'application/json' }};
    let res_post = http.post(`${BASE_URL}/api/cart/editCart`, payload_post, params_post);
  
    const checkRes = check(res_post, {
      'status is 200': (r) => r.status === 200,
    });  

    if(res_post.status == 200) {
      let data = JSON.parse(res_post.body).data;
      cart[cartid] = JSON.parse(JSON.stringify(data["carts"]));
      client.set('cart_data ' + __VU, JSON.stringify(data["carts"]), 0);
    }  
  
    sleep(0.5);
}


export function PurchaseFromCarts() {

    let cart = JSON.parse(client.get('cart_data ' + __VU));
    let cartids = [];

    for (let items = 0; items < cart.length; items ++)
      cartids.push(cart[items]["ID"]);
  
    const payload_post = JSON.stringify({ 'accountID': __VU, 'cartIDs': cartids });
    const params_post = { headers: { 'Content-Type': 'application/json' } };
    let res_post = http.post(`${BASE_URL}/api/purchase/sync`, payload_post, params_post);  
  
    const checkRes = check(res_post, {
      'status is 200': (r) => r.status === 200,
    });
  
    sleep(0.5);
}
