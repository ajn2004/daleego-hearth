package com.daleego.hearth.api

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import retrofit2.HttpException

data class ApiError(
    val error: String
)

private val moshi = Moshi.Builder()
    .add(KotlinJsonAdapterFactory())
    .build()
private val apiErrorAdapter = moshi.adapter(ApiError::class.java)

fun readableApiError(throwable: Throwable): String {
    return try {
        if (throwable is HttpException) {
            val errorBody = throwable.response()?.errorBody()?.string()

            if (!errorBody.isNullOrBlank()) {
                return try {
                    apiErrorAdapter.fromJson(errorBody)?.error
                        ?: "HTTP ${throwable.code()}: ${throwable.message()}"
                } catch (_: Exception) {
                    errorBody
                }
            }

            return "HTTP ${throwable.code()}: ${throwable.message()}"
        }

    throwable.message ?: "Failed to connect to Hearth."
} catch (_: Exception) {
    "Failed to pair device."
}
}