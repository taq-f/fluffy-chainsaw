<template>
  <div class="path">
    <div class="dir">
      {{ dir }}
    </div>
    <div class="sep">
      /
    </div>
    <div class="filename">
      {{ filename }}
    </div>
  </div>
</template>

<script>
export default {
  props: {
    path: String
  },
  data() {
    return {
      dir: "",
      filename: ""
    };
  },
  created() {
    const lastIndex = this.path.lastIndexOf("/");
    this.dir = this.path.substring(0, lastIndex);
    this.filename = this.path.substring(lastIndex + 1);
  }
};
</script>


<style lang="scss" scoped>
.path {
  display: flex;
  white-space: nowrap;
}

.dir {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  color: #607d8b;
}

.sep {
  color: #607d8b;
}

.filename {
  flex: 1;
  color: #263238;
}
</style>
