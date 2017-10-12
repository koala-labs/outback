<template>
  <div id="app">
    <router-view></router-view>
    <notifications group="ufo"
                   classes="n-light"
                   :max="3"
                   :width="350"/>
  </div>
</template>

<script>
export default {
  name: 'app',
};
</script>

<style lang='scss'>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 10px;
}

.notification.n-light {
  margin: 0.5rem 0.5rem 1rem;
  border-radius: 2px;

  font-size: 14px;
  font-weight: bold;

  padding: 10px 20px;

  color: #03a9f4;
  background: #E1F5FE;

  transition: box-shadow .25s;
  
  box-shadow: 0 2px 5px 0 rgba(0,0,0,0.16), 0 2px 10px 0 rgba(0,0,0,0.12);

  &.error {
    background: #FFEBEE;
    color: #F44336;

    .notification-title {
      color: #F44336;
    }
  }

  &.success {
    background: #E8F5E9;
    color: #4CAF50;

    .notification-title {
      color: #4CAF50;
    }
  }

  .notification-title {
    letter-spacing: 1px;
    text-transform: uppercase;
    font-size: 14px;
    color: #03a9f4;
  }
}
</style>
