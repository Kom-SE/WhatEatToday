<template>
  <div class="container mx-auto p-4">
    <h1 class="text-2xl font-bold mb-4">热门美食推荐</h1>

    <!-- 菜谱网格 -->
    <div class="grid grid-cols-3 gap-6">
      <div
        v-for="recipe in paginatedRecipes"
        :key="recipe.id"
        class="bg-gray-800 rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105"
      >
        <img
          :src="recipe.images"
          alt="Recipe Image"
          class="w-full h-48 object-cover"
        />
        <div class="p-4">
          <h2 class="text-xl font-semibold text-white">{{ recipe.name }}</h2>
          <p class="text-gray-300 text-sm mt-1">{{ recipe.description }}</p>
          <div class="flex justify-between items-center mt-3 text-gray-400">
            <span>麻辣度：无</span>
            <span>👍 {{ recipe.likes }} 好评</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 分页控件 -->
    <div class="mt-6 flex justify-center items-center space-x-2">
      <button
        @click="prevPage"
        :disabled="currentPage === 1"
        class="px-3 py-1 bg-gray-700 text-white rounded disabled:opacity-50"
      >
        上一页
      </button>

      <span class="text-gray-600">
        第 {{ currentPage }} 页 / 共 {{ totalPages }} 页
      </span>

      <button
        @click="nextPage"
        :disabled="currentPage === totalPages"
        class="px-3 py-1 bg-gray-700 text-white rounded disabled:opacity-50"
      >
        下一页
      </button>
    </div>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  data() {
    return {
      recipes: [],         // 所有菜谱
      currentPage: 1,      // 当前页码
      pageSize: 6,         // 每页显示条数
    };
  },
  mounted() {
    this.fetchRecipes();
  },
  methods: {
    async fetchRecipes() {
      try {
        const res = await axios.get('http://127.0.0.1:3000/recipe/getall');
        this.recipes = res.data.recipes;
      } catch (error) {
        console.error('请求失败:', error);
      }
    },
    nextPage() {
      if (this.currentPage < this.totalPages) {
        this.currentPage++;
      }
    },
    prevPage() {
      if (this.currentPage > 1) {
        this.currentPage--;
      }
    },
  },
  computed: {
    totalPages() {
      return Math.ceil(this.recipes.length / this.pageSize);
    },
    paginatedRecipes() {
      const start = (this.currentPage - 1) * this.pageSize;
      const end = start + this.pageSize;
      return this.recipes.slice(start, end);
    },
  },
};
</script>

<style scoped>
/* 可选样式优化 */
.container {
  max-width: 1200px;
}
</style>