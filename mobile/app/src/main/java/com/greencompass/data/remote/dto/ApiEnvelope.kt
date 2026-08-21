package com.greencompass.data.remote.dto

import com.google.gson.annotations.SerializedName

data class ApiResponse<T>(
    @SerializedName("data") val data: T?,
    @SerializedName("meta") val meta: ApiMeta? = null,
    @SerializedName("error") val error: ApiError? = null
)

data class ApiMeta(
    @SerializedName("request_id") val requestId: String? = null,
    @SerializedName("generated_at") val generatedAt: String? = null
)

data class ApiError(
    @SerializedName("code") val code: String,
    @SerializedName("message") val message: String,
    @SerializedName("request_id") val requestId: String? = null
)
