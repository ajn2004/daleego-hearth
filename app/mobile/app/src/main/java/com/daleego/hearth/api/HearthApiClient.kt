package com.daleego.hearth.api

import com.squareup.moshi.Moshi
import com.squareup.moshi.kotlin.reflect.KotlinJsonAdapterFactory
import okhttp3.OkHttpClient
import retrofit2.Retrofit
import retrofit2.converter.moshi.MoshiConverterFactory

object HearthApiClient {
    fun create(
        baseUrl: String,
        getDeviceApiKey: () -> String?
    ): HearthApi {
        val moshi = Moshi.Builder()
            .add(KotlinJsonAdapterFactory())
            .build()

        val client = OkHttpClient.Builder()
            .addInterceptor { chain ->
                val apiKey = getDeviceApiKey()
                android.util.Log.d(
                    "HearthApiClient",
                    "Device API key present: ${!apiKey.isNullOrBlank()}"
                )
                val requestBuilder = chain.request().newBuilder()
                if (!apiKey.isNullOrBlank()) {
                    requestBuilder.addHeader("X-Device-Api-Key", apiKey)
                }
                chain.proceed(requestBuilder.build())
            }
            .build()

        return Retrofit.Builder()
            .baseUrl(baseUrl)
            .client(client)
            .addConverterFactory(MoshiConverterFactory.create(moshi))
            .build()
            .create(HearthApi::class.java)
    }
}
