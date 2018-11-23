<template>
  <div>
    <div class="todo" v-for="todo in todos" :key="todo.id">
      <div class="todo-element todo-file">
        <filepath :path="todo.file"></filepath>
      </div>
      <div class="todo-element todo-line">
        {{ todo.line }}
      </div>
      <div class="todo-element todo-text">
        {{ todo.text }}
      </div>
    </div>
  </div>
</template>

<script>
import Filepath from "@/components/Filepath.vue";
import axios from "axios";

// const dummyTodos = [
//   {
//     id: "todo1",
//     file: "/src/sample.js",
//     line: 12,
//     text: "// TODO あれをこうしてああする。todo1"
//   },
//   {
//     id: "todo3",
//     file: "/src/sample.js",
//     line: 1,
//     text: "// TODO あれをこうしてああする。todo3"
//   },
//   {
//     id: "todo2",
//     file: "/src/sample.js",
//     line: 66,
//     text: "// TODO あれをこうしてああする。todo2"
//   },
//   {
//     id: "todo4",
//     file: "/src/sample.ts",
//     line: 3,
//     text: "// TODO あれをこうしてああする。todo4"
//   },
//   {
//     id: "todo5",
//     file: "/src/sample.ts",
//     line: 99,
//     text: "// TODO あれをこうしてああする。todo5"
//   }
// ];

export default {
  components: {
    Filepath
  },
  data() {
    return {
      todos: []
    };
  },
  async mounted() {
    const response = await axios.get("/api/todos");
    const todos = response.data;

    this.todos = todos.sort((todo1, todo2) => {
      if (todo1.file !== todo2.file) {
        return todo1.file > todo2.file ? 1 : -1;
      } else {
        return todo1.line - todo2.line;
      }
    });
  }
};
</script>

<style lang="scss" scoped>
.todo {
  display: flex;
  border-bottom: 1px solid #cfd8dc;
  &:hover {
    background-color: #cfd8dc;
  }

  .todo-element {
    margin-top: 15px;
    margin-bottom: 15px;
    margin-left: 20px;
    &:last-of-type {
      margin-right: 20px;
    }

    &.todo-file {
      width: 150px;
      overflow: hidden;
    }
    &.todo-line {
      width: 40px;
      text-align: right;
      overflow: hidden;
    }
    &.todo-text {
      flex: 1;
    }
  }
}
</style>
