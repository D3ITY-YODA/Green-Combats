package com.greencompass.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    @Serializable data object Welcome : AppRoute
    @Serializable data object LanguageSelection : AppRoute
    @Serializable data object AccountChoice : AppRoute
    @Serializable data object SignIn : AppRoute
    @Serializable data object SignUp : AppRoute
    @Serializable data object PlaceSetup : AppRoute
    @Serializable data object Interests : AppRoute
    @Serializable data object Permission : AppRoute
    @Serializable data object SetupComplete : AppRoute

    @Serializable data object Today : AppRoute
    @Serializable data object Explore : AppRoute
    @Serializable data object Updates : AppRoute
    @Serializable data object Report : AppRoute
    @Serializable data object Profile : AppRoute

    @Serializable data class ReportForm(val reportType: String) : AppRoute
    @Serializable data object ReportSuccess : AppRoute
    @Serializable data object UpdateDetail : AppRoute
    @Serializable data object SavedPlaces : AppRoute
    @Serializable data object LocationSwitcher : AppRoute
    @Serializable data object Settings : AppRoute
    @Serializable data object Notifications : AppRoute
    @Serializable data object Help : AppRoute
    @Serializable data object About : AppRoute
    @Serializable data object Offline : AppRoute
}
