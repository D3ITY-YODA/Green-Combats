package com.greencompass.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    // Onboarding
    @Serializable data object Welcome : AppRoute
    @Serializable data object LanguageSelection : AppRoute
    @Serializable data object AccountChoice : AppRoute
    @Serializable data object PlaceSetup : AppRoute
    @Serializable data object Interests : AppRoute
    @Serializable data object SetupComplete : AppRoute

    // Main App
    @Serializable data object Today : AppRoute
    @Serializable data object Explore : AppRoute
    @Serializable data object Updates : AppRoute
    @Serializable data object Report : AppRoute
    @Serializable data object Profile : AppRoute
}
