package com.greencompass.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    @Serializable
    data object Welcome : AppRoute

    @Serializable
    data object LanguageSelection : AppRoute

    @Serializable
    data object AccountChoice : AppRoute
}
