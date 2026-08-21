package com.greencompass.feature.reports

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.greencompass.core.ui.GreenCompassColors

@Composable
fun ReportFormScreen(
    reportType: String,
    onBack: () -> Unit,
    onSubmitSuccess: () -> Unit
) {
    var description by remember { mutableStateOf("") }
    var isSubmitting by remember { mutableStateOf(false) }

    Surface(modifier = Modifier.fillMaxSize(), color = GreenCompassColors.WarmWhite) {
        Column(modifier = Modifier.fillMaxSize().padding(24.dp)) {
            Row(verticalAlignment = androidx.compose.ui.Alignment.CenterVertically) {
                IconButton(onClick = onBack) { Icon(Icons.Default.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal) }
                Text(text = reportType, fontSize = 20.sp, fontWeight = FontWeight.SemiBold, color = GreenCompassColors.Charcoal)
            }
            Spacer(Modifier.height(24.dp))
            
            Text(text = "Tell us more about what you are seeing.", fontSize = 15.sp, color = GreenCompassColors.MutedText)
            Spacer(Modifier.height(16.dp))
            
            OutlinedTextField(
                value = description,
                onValueChange = { description = it },
                label = { Text("Description") },
                placeholder = { Text("Optional details...") },
                modifier = Modifier.fillMaxWidth().height(150.dp),
                shape = RoundedCornerShape(12.dp)
            )
            
            Spacer(Modifier.weight(1f))
            
            Button(
                onClick = { 
                    isSubmitting = true
                    // TODO: Call ViewModel to submit to Room DB
                    onSubmitSuccess() 
                },
                modifier = Modifier.fillMaxWidth().height(56.dp),
                enabled = !isSubmitting,
                colors = ButtonDefaults.buttonColors(containerColor = GreenCompassColors.ForestGreen),
                shape = RoundedCornerShape(12.dp)
            ) {
                Text(text = if (isSubmitting) "Sending..." else "Send update", fontSize = 16.sp)
            }
        }
    }
}
