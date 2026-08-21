package com.greencompass.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    // Onboarding
    @Serializable data object Welcome : AppRoute
    @Serializable data object LanguageSelection : AppRoute
    @Serializable data object AccountChoice : AppRoute
    @Serializable data object PersonalRegistration : AppRoute
    @Serializable data object GoogleSignIn : AppRoute
    @Serializable data object PhoneEmailSignIn : AppRoute
    @Serializable data object VerificationCode : AppRoute
    @Serializable data object OrganizationSearch : AppRoute
    @Serializable data object OrganizationSelection : AppRoute
    @Serializable data object RequestAccess : AppRoute
    @Serializable data object InvitationAcceptance : AppRoute
    @Serializable data object AccessPending : AppRoute
    @Serializable data object LocationSelection : AppRoute
    @Serializable data object SearchPlace : AppRoute
    @Serializable data object MapPlaceSelection : AppRoute
    @Serializable data object Interests : AppRoute
    @Serializable data object NotificationPreferences : AppRoute
    @Serializable data object PrivacyPermission : AppRoute
    @Serializable data object SetupComplete : AppRoute
    
    // Main App
    @Serializable data object Today : AppRoute
    @Serializable data object Explore : AppRoute
    @Serializable data object LocalOutlook : AppRoute
    @Serializable data object SeasonalInformation : AppRoute
    @Serializable data object WaterOutlook : AppRoute
    @Serializable data object LandEcosystems : AppRoute
    @Serializable data object FoodAgriculture : AppRoute
    @Serializable data object CommunityUpdates : AppRoute
}
