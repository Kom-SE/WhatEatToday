<template>
  <div class="container">
    <div class="page-header">
      <p class="page-title">发布菜谱</p>
      <p class="page-subtitle">分享您的美食创意</p>
    </div>
    <div class="mainpage">
      <div class="leftpage">
        <div class="form-section">
          <label for="recipe-name" class="section-label">菜谱名称</label>
          <el-input
            id="recipe-name"
            v-model="recipe.name"
            placeholder="请输入菜谱名称"
            size="large"
          />
        </div>

        <div class="form-section">
          <label for="recipe-description" class="section-label">描述</label>
          <el-input
            id="recipe-description"
            v-model="recipe.description"
            type="textarea"
            :rows="3"
            placeholder="请简要描述您的菜谱，例如口味特点、适合人群等"
            size="large"
          />
        </div>

        <div class="form-section">
          <label for="recipe-steps" class="section-label">制作过程</label>
          <el-input
            id="recipe-steps"
            v-model="recipe.steps"
            type="textarea"
            :rows="6"
            placeholder="请详细描述制作步骤，每一步骤用回车分隔"
            size="large"
          />
        </div>

        <div class="form-section">
          <label class="section-label">所需食材</label>

          <!-- 已选择的食材显示区域 -->
          <div
            class="selected-ingredients"
            v-if="selectedIngredients.length > 0"
          >
            <el-tag
              v-for="ingredient in selectedIngredients"
              :key="ingredient.id"
              closable
              @close="removeIngredient(ingredient.id)"
              size="large"
              class="ingredient-tag"
            >
              {{ ingredient.name }}
            </el-tag>
          </div>

          <!-- 食材搜索输入框 -->
          <el-input
            v-model="searchKeyword"
            placeholder="输入食材名称进行搜索..."
            @input="searchIngredients"
            @keyup.enter="handleEnterKey"
            size="large"
            :loading="isSearching"
            clearable
          >
            <template #suffix>
              <el-icon v-if="isSearching"><Loading /></el-icon>
              <el-icon v-else><Search /></el-icon>
            </template>
          </el-input>

          <!-- 搜索结果下拉列表 -->
          <div v-if="searchResults.length > 0" class="search-results">
            <div
              v-for="result in searchResults"
              :key="result.id"
              class="search-result-item"
              @click="selectIngredient(result)"
            >
              {{ result.name }}
            </div>
          </div>

          <!-- 创建新食材对话框 -->
          <el-dialog
            v-model="showCreateDialog"
            title="创建新食材"
            width="400px"
          >
            <el-form :model="newIngredientForm" label-width="80px">
              <el-form-item label="食材名称">
                <el-input
                  v-model="newIngredientForm.name"
                  placeholder="请输入食材名称"
                />
              </el-form-item>

              <el-form-item label="食材描述">
                <el-input
                  v-model="newIngredientForm.description"
                  type="textarea"
                  :rows="3"
                  placeholder="请输入食材描述"
                />
              </el-form-item>

              <el-form-item label="食材图片">
                <el-upload
                  class="upload-demo"
                  :before-upload="beforeUpload"
                  :on-change="handleImageChange"
                  :on-remove="handleImageRemove"
                  :file-list="fileList"
                  list-type="picture"
                  :auto-upload="false"
                  accept="image/*"
                >
                  <el-button type="primary">选择图片</el-button>
                  <template #tip>
                    <div class="el-upload__tip">
                      只能上传jpg/png文件，且不超过5MB
                    </div>
                  </template>
                </el-upload>
              </el-form-item>
            </el-form>

            <template #footer>
              <el-button @click="showCreateDialog = false">取消</el-button>
              <el-button type="primary" @click="createNewIngredient"
                >创建</el-button
              >
            </template>
          </el-dialog>
        </div>

        <div class="form-section">
          <label for="prep-time" class="section-label">制作时间</label>
          <el-time-picker
            id="prep-time"
            v-model="recipe.prepTime"
            placeholder="请选择制作时间"
            format="HH:mm:ss"
            value-format="HH:mm:ss"
            size="large"
            class="time-select"
            clearable
          />
        </div>

        <div class="form-section">
          <label for="flavors" class="section-label">风味特色</label>
          <el-select
            id="flavors"
            v-model="recipe.flavors"
            multiple
            collapse-tags
            placeholder="选择菜品的风味特色"
            clearable
            size="large"
            class="flavor-select"
          >
            <el-option
              v-for="flavor in availableFlavors"
              :key="flavor.value"
              :label="flavor.label"
              :value="flavor.value"
            />
          </el-select>
        </div>

        <el-button
          type="primary"
          size="large"
          class="submit-button"
          @click="submitRecipe"
          >发布菜谱</el-button
        >
      </div>

      <div class="rightpage">
        <h3>菜谱预览</h3>
        <p><strong>名称:</strong> {{ recipe.name || "未填写" }}</p>
        <p><strong>描述:</strong> {{ recipe.description || "未填写" }}</p>
        <p><strong>制作时间:</strong> {{ recipe.prepTime || "未选择" }}</p>
        <!-- 显示实际选择的食材 -->
        <p>
          <strong>已选食材:</strong>
          <span v-if="selectedIngredients.length > 0">
            {{ selectedIngredients.map((item) => item.name).join(", ") }}
          </span>
          <span v-else>未选择</span>
        </p>
        <!-- 显示风味标签而不是value -->
        <p>
          <strong>风味特色:</strong>
          <span v-if="recipe.flavors.length > 0">
            {{ getFlavorLabels(recipe.flavors).join(", ") }}
          </span>
          <span v-else>未选择</span>
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { ElMessage } from "element-plus";
import { Search, Loading } from "@element-plus/icons-vue";
import type { UploadProps, UploadFile } from "element-plus";

interface Ingredient {
  id: number;
  name: string;
  description?: string; // 添加description
  image?: string; // 添加image
  category?: string;
}

interface NewIngredientForm {
  name: string;
  description: string;
  category: string;
}

// 已选择的食材列表
const selectedIngredients = ref<Ingredient[]>([]);

// 搜索相关
const searchKeyword = ref("");
const searchResults = ref<Ingredient[]>([]);
const isSearching = ref(false);

// 创建新食材相关
const showCreateDialog = ref(false);
const newIngredientForm = reactive<NewIngredientForm>({
  name: "",
  description: "",
  category: "",
});

// 防抖搜索函数
let searchTimer: NodeJS.Timeout | null = null;

const searchIngredients = async () => {
  if (!searchKeyword.value.trim()) {
    searchResults.value = [];
    return;
  }

  // 清除之前的定时器
  if (searchTimer) {
    clearTimeout(searchTimer);
  }

  // 设置新的定时器，500ms后执行搜索
  searchTimer = setTimeout(async () => {
    isSearching.value = true;

    try {
      const response = await fetch(
        `http://localhost:3000/food/get/${encodeURIComponent(
          searchKeyword.value.trim()
        )}`,
        {
          method: "GET",
          headers: {
            Authorization: localStorage.getItem("jwt") || "",
            "Content-Type": "application/json",
          },
        }
      );

      if (response.ok) {
        const data = await response.json();

        // 修复：检查正确的数据结构
        if (data.food && data.food.name) {
          // 后端返回了food对象，转换为数组格式
          const foodItem: Ingredient = {
            id: Date.now(), // 临时生成ID，如果后端没有提供
            name: data.food.name,
            description: data.food.description,
            image: data.food.image,
          };
          searchResults.value = [foodItem];
        } else {
          searchResults.value = [];
        }
      }
    } catch (error) {
      console.error("搜索食材失败:", error);
      searchResults.value = [];
    } finally {
      isSearching.value = false;
    }
  }, 500);
};

// 处理回车键
const handleEnterKey = () => {
  if (searchResults.value.length === 0 && searchKeyword.value.trim()) {
    // 没有搜索结果，提示创建新食材
    newIngredientForm.name = searchKeyword.value.trim();
    showCreateDialog.value = true;
  } else if (searchResults.value.length > 0) {
    // 有搜索结果，选择第一个
    selectIngredient(searchResults.value[0]);
  }
};

// 选择食材
const selectIngredient = (ingredient: Ingredient) => {
  // 检查是否已经选择过
  const isAlreadySelected = selectedIngredients.value.some(
    (item) => item.id === ingredient.id
  );

  if (isAlreadySelected) {
    ElMessage.warning("该食材已经添加过了");
    return;
  }

  // 添加到已选择列表
  selectedIngredients.value.push(ingredient);

  // 清空搜索
  searchKeyword.value = "";
  searchResults.value = [];

  ElMessage.success(`已添加食材: ${ingredient.name}`);
};

// 移除食材
const removeIngredient = (ingredientId: number) => {
  selectedIngredients.value = selectedIngredients.value.filter(
    (item) => item.id !== ingredientId
  );
  ElMessage.success("已移除食材");
};

// 图片上传相关
const fileList = ref<UploadFile[]>([]);
const selectedImageFile = ref<File | null>(null);

// 图片上传处理函数
const beforeUpload: UploadProps["beforeUpload"] = (file) => {
  const isImage = file.type.startsWith("image/");
  const isLt2M = file.size / 1024 / 1024 < 2;

  if (!isImage) {
    ElMessage.error("只能上传图片文件!");
    return false;
  }
  if (!isLt2M) {
    ElMessage.error("图片大小不能超过 2MB!");
    return false;
  }
  return false; // 阻止自动上传
};

const handleImageChange: UploadProps["onChange"] = (file) => {
  if (file.raw) {
    selectedImageFile.value = file.raw;
  }
};

const handleImageRemove = () => {
  selectedImageFile.value = null;
  fileList.value = [];
};

// 取消创建
const cancelCreate = () => {
  showCreateDialog.value = false;
  newIngredientForm.name = "";
  newIngredientForm.description = "";
  newIngredientForm.category = "";
  selectedImageFile.value = null;
  fileList.value = [];
};

// 创建新食材 - 使用FormData
const createNewIngredient = async () => {
  if (!newIngredientForm.name.trim()) {
    ElMessage.error("请输入食材名称");
    return;
  }

  try {
    // 创建FormData对象
    const formData = new FormData();
    formData.append("name", newIngredientForm.name.trim());
    formData.append("description", newIngredientForm.description.trim());

    // 如果选择了图片，添加到FormData
    if (selectedImageFile.value) {
      formData.append("image", selectedImageFile.value);
    }

    const response = await fetch("http://localhost:3000/food/add", {
      method: "POST",
      headers: {
        Authorization: localStorage.getItem("jwt") || "",
        // 注意：使用FormData时不要设置Content-Type，让浏览器自动设置
      },
      body: formData,
    });

    if (response.ok) {
      const data = await response.json();
      const newIngredient: Ingredient = {
        id: data.food?.id || Date.now(),
        name: data.food?.name || newIngredientForm.name,
        description: data.food?.description || newIngredientForm.description,
        image: data.food?.image || "",
        category: newIngredientForm.category,
      };

      // 添加到已选择列表
      selectedIngredients.value.push(newIngredient);

      // 重置表单
      cancelCreate();

      ElMessage.success(`成功创建并添加食材: ${newIngredient.name}`);
    } else {
      const errorData = await response.json();
      ElMessage.error(errorData.message || "创建食材失败");
    }
  } catch (error) {
    console.error("创建食材失败:", error);
    ElMessage.error("创建食材失败");
  }
};

// 暴露给父组件使用的数据
defineExpose({
  selectedIngredients,
});

interface Recipe {
  name: string;
  description: string;
  steps: string;
  ingredients: string[];
  prepTime: string;
  flavors: string[];
}

const recipe = reactive<Recipe>({
  name: "",
  description: "",
  steps: "",
  ingredients: [],
  prepTime: "",
  flavors: [],
});

const newIngredient = ref("");

const addCustomIngredient = () => {
  if (
    newIngredient.value &&
    !recipe.ingredients.includes(newIngredient.value)
  ) {
    recipe.ingredients.push(newIngredient.value);
    availableIngredients.value.push({
      label: newIngredient.value,
      value: newIngredient.value,
    }); // 也可以添加到可用列表中
    newIngredient.value = "";
    ElMessage.success("自定义食材已添加");
  } else if (
    newIngredient.value &&
    recipe.ingredients.includes(newIngredient.value)
  ) {
    ElMessage.warning("该食材已在列表中");
  }
};

// 添加获取风味标签的函数
const getFlavorLabels = (flavorValues: string[]) => {
  return flavorValues.map((value) => {
    const flavor = availableFlavors.value.find((f) => f.value === value);
    return flavor ? flavor.label : value;
  });
};
const availableFlavors = ref([
  { label: "麻辣", value: "spicy_numbing" },
  { label: "酸甜", value: "sweet_sour" },
  { label: "咸香", value: "savory" },
  { label: "清淡", value: "light" },
  { label: "香辣", value: "spicy" },
  { label: "鲜美", value: "umami" },
]);

const submitRecipe = () => {
  console.log("提交的菜谱数据:", recipe);
  ElMessage.success("菜谱发布成功！ (数据已打印到控制台)");
  // 在这里可以添加实际的API调用逻辑
};
</script>

<style scoped>
/* 页面标题样式 - 添加这部分 */
.page-header {
  text-align: center;
  margin-bottom: 20px;
}

.page-title {
  font-size: 30px; /* 调整标题字体大小 */
  font-weight: 600;
  color: #303133;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #409eff, #67c23a);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  font-size: 14px; /* 调整副标题字体大小 */
  color: #909399;
  margin: 0;
  font-weight: 400;
}

/* 字体和全局间距 */
.container {
  font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
    "Microsoft YaHei", sans-serif;
  padding: 0px 10px;
  background-color: #f4f6f9;
  height: 100vh; /* 设置容器高度 */
  overflow: hidden; /* 防止整体滚动 */
}

/* 主布局 */
.mainpage {
  display: flex;
  justify-content: center;
  align-items: flex-start;
  margin-top: 30px;
  gap: 30px;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
  height: calc(100vh - 100px); /* 设置主页面高度 */
}

.leftpage,
.rightpage {
  flex: 1;
  min-width: 400px;
  background-color: #ffffff;
  border-radius: 12px;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.08);
  padding: 30px;
  display: flex;
  flex-direction: column;
  gap: 25px;

  /* 关键修改：设置高度和滚动 */
  height: 100%; /* 占满父容器高度 */
  overflow-y: auto; /* 两边都可以滚动 */
  max-height: calc(100vh - 140px); /* 限制最大高度 */
}

.leftpage {
  border-left: 5px solid #409eff;
}

.rightpage {
  border-right: 5px solid #67c23a;
  /* 移除之前的sticky定位，让两边保持一致 */
}

/* 表单分段样式 */
.form-section {
  display: flex;
  flex-direction: column;
  gap: 10px;
  flex-shrink: 0; /* 防止内容被压缩 */
}

.section-label {
  font-size: 16px;
  color: #303133;
  font-weight: 600;
  margin-bottom: 5px;
  display: flex;
  align-items: center;
}

/* 已选择的食材显示区域 */
.selected-ingredients {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 10px;
  background-color: #f5f7fa;
  border-radius: 6px;
  min-height: 50px;
  align-items: center;
  margin-bottom: 10px;
}

.ingredient-tag {
  margin: 2px;
}

/* 搜索结果下拉列表 */
.search-results {
  border: 1px solid #dcdfe6;
  border-radius: 4px;
  background-color: white;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  max-height: 200px;
  overflow-y: auto;
  z-index: 100;
  margin-top: 5px;
}

.search-result-item {
  padding: 10px 15px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background-color 0.2s;
}

.search-result-item:hover {
  background-color: #f5f7fa;
}

.search-result-item:last-child {
  border-bottom: none;
}

/* El-input, El-textarea 样式 */
.el-input,
.el-textarea {
  width: 100%;
}

/* 食材多选框组 */
.ingredients-checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.ingredients-checkbox-group .el-checkbox.is-bordered {
  border-radius: 6px;
  padding: 10px 18px;
  transition: all 0.2s ease;
}

.ingredients-checkbox-group .el-checkbox.is-bordered.is-checked {
  background-color: #ecf5ff;
  border-color: #409eff;
  color: #409eff;
}

/* 自定义食材输入框间距 */
.mt-3 {
  margin-top: 15px;
}

/* 时间选择器 */
.time-select {
  width: 100%;
}

/* 风味选择器 */
.flavor-select {
  width: 100%;
}

/* 提交按钮 */
.submit-button {
  margin-top: 20px;
  align-self: flex-end;
  padding: 12px 30px;
  font-size: 16px;
  border-radius: 8px;
  flex-shrink: 0; /* 防止按钮被压缩 */
}

/* 右侧预览样式 */
.rightpage h3 {
  color: #409eff;
  margin-bottom: 20px;
  text-align: center;
  font-size: 22px;
  padding-bottom: 10px;
  border-bottom: 2px solid #e0e0e0;
  flex-shrink: 0; /* 标题不被压缩 */
}

.rightpage p {
  line-height: 1.8;
  margin-bottom: 10px;
  color: #555;
  font-size: 15px;
  flex-shrink: 0; /* 内容不被压缩 */
}

.rightpage strong {
  color: #333;
}

/* 图片上传样式 */
.upload-demo {
  width: 100%;
}

.el-upload__tip {
  color: #999;
  font-size: 12px;
  margin-top: 5px;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .container {
    height: auto; /* 移动端取消固定高度 */
    overflow: visible;
  }

  .mainpage {
    flex-direction: column;
    margin-top: 20px;
    gap: 20px;
    height: auto; /* 移动端取消固定高度 */
  }

  .leftpage,
  .rightpage {
    min-width: unset;
    width: 100%;
    padding: 20px;
    height: auto; /* 移动端取消固定高度 */
    max-height: none; /* 移动端取消最大高度限制 */
    overflow-y: visible; /* 移动端取消滚动 */
  }

  .submit-button {
    align-self: stretch;
  }
  .preview-section {
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 1px solid #f0f0f0;
  }

  .preview-section:last-child {
    border-bottom: none;
  }

  .preview-content {
    margin-top: 8px;
    padding: 10px;
    background-color: #f8f9fa;
    border-radius: 6px;
    font-size: 14px;
    line-height: 1.6;
  }

  .steps-content .step-item {
    margin-bottom: 8px;
    padding: 5px 0;
  }

  .ingredients-list,
  .flavors-list {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
  }

  .ingredient-preview-tag,
  .flavor-preview-tag {
    margin: 2px;
  }

  .placeholder {
    color: #999;
    font-style: italic;
  }
}
</style>
